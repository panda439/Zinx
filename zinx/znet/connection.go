package znet

import (
	"net"
	"zinx/ziface"
	"fmt"
	"io"
	"errors"
)

//具体的TCP链接模块
type Connection struct {
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


}

func NewConnection(conn *net.TCPConn,connID uint32,handler ziface.IMsgHandle) ziface.IConnection {

	return &Connection{
		Conn:conn,
		ConnID:connID,
		MsgHandler:handler,
		IsClosed:false,
	}
}

//针对链接读业务的方法
func (c *Connection) StartReader() {
	//从对端读数据
	fmt.Println("Reader go is startin....")
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
		go c.MsgHandler.DoMsgHandler(req)




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




}
//停止链接
func (c *Connection)Stop(){
	fmt.Println("c.Stop()...ConnId =",c.ConnID)
	//回收工作
	if c.IsClosed == true {
		return
	}

	//关闭原生套接字
	_ = c.Conn.Close()
	c.IsClosed = true






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


	if _,err := c.Conn.Write(binaryMsg);err!=nil{
		fmt.Println("send buf error")
		return err
	}
	return nil
}