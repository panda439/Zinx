package znet

import (
	"net"
	"zinx/ziface"
	"fmt"
	"io"
	"errors"
	"zinx/config"
)

//具体的TCP链接模块
type Connection struct {
	//当前链接属于哪个sever
	server ziface.IServer

	//当前链接的原生套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前的链接状态
	IsClosed bool

	//当前链接所绑定的业务处理方法
	//HandleAPI ziface.HandleFunc
	//当前链接所绑定的router
	MsgHandler ziface.IMsgHandle


	//添加一个Reader和Writer通信的Channel
	msgChan chan []byte

	//创建一个Channel 用来Reader通知Writer conn已经关闭，需要推出的消息
	writerExitChan chan bool



}

func NewConnection(server ziface.IServer,conn *net.TCPConn,connID uint32,handler ziface.IMsgHandle) ziface.IConnection {

	 c :=  &Connection{
		server:server,
		Conn:conn,
		ConnID:connID,
		MsgHandler:handler,
		IsClosed:false,
		msgChan:make(chan []byte),
		writerExitChan:make(chan bool),
	}
	//当已经成功创建一个链接的时候，添加到链接管理器中
	c.server.GetConnMgr().Add(c)

	return c



}
func (c *Connection) StartWriter() {
	fmt.Println("【Writer Goroutine is Started】")
	defer fmt.Println("【Writer Goroutine is Stop】")

	//IO多路复用
	for {
		select {
		case data :=<-c.msgChan:
			//有数据需要写给客户端
			if _ ,err := c.Conn.Write(data);err!=nil{
				fmt.Println("Send data error ",err)
				return
			}
		case <-c.writerExitChan:
			return
		}


	}

}





//针对链接读业务的方法
func (c *Connection) StartReader() {
	//从对端读数据
	fmt.Println("【Reader go is startin....】")
	defer 	fmt.Println("【Reader go is stop....】")

	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is = ", c.GetRemoteAddr().String())

	defer c.Stop()

	for  {
		//buf := make([]byte,config.GlobalObject.MaxPackageSize)
		//cnt,err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read",err)
		//	break
		//}
		dp := NewDataPack()

		headData:= make([]byte,dp.GetHeadLen())

		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("read msg head error",err)
			break
		}


		msg,err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Unpack error",err)
			break
		}

		//根据长度再次读取
		var data []byte
		if msg.GetMsgLen() > 0 {
			//有内容
			data = make([]byte,msg.GetMsgLen())
			if _,err := io.ReadFull(c.Conn,data);err!=nil{
				fmt.Println(err)
				break
			}
		}

		msg.SetData(data)



		//将当前一次性
		req := NewRequest(c,msg)

		//调用用户传递进来的业务 模版 设计模式
		//将req交给worker工作池来处理
		if config.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}





		//将数据 传递给我们 定义好的Handle Callback方法
		/*if err := c.HandleAPI(req); err != nil {
			fmt.Println("ConnID", c.ConnID, "Handle is error", err)
			break
		}*/

	}


}



//启动链接
func (c *Connection)Start(){
	fmt.Println("Conn Start()...id=",c.ConnID)


	//先进行读业务
	go c.StartReader()
	//进行写业务
	go c.StartWriter()

	c.server.CallOnConnStart(c)


}
//停止链接
func (c *Connection)Stop(){
	fmt.Println("c.Stop()...ConnId =",c.ConnID)

	c.server.CallOnConnStop(c)
	//回收工作
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	c.writerExitChan<-true

	//关闭原生套接字
	_ = c.Conn.Close()

	//
	c.server.GetConnMgr().Remove(c.ConnID)


	close(c.writerExitChan)
	close(c.msgChan)

}
//获取链接ID
func (c *Connection)GetConnID() uint32{
	return c.ConnID

}
//获取conn的原生socket套接字
func (c *Connection)GetTCPConnection() *net.TCPConn{
	return c.Conn

}

//获取远程客户端的ip地址
func (c *Connection)GetRemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}

//发送数据给对方客户端
func (c *Connection)Send(msgId uint32,msgData []byte ) error{
	if c.IsClosed == true {
		return errors.New("connection closed..send Msg")
	}
	//封装成msg
	dp := NewDataPack()
	binaryMsg,err:=dp.Pack(NewMsgPackage(msgId,msgData))
	if err != nil {
		fmt.Println(err)
		return err
	}

	//将要发送的打包好的二进制数发送channel让write
	c.msgChan <- binaryMsg

	return nil
}
