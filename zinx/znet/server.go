package znet

import (
	"zinx/ziface"
	"fmt"
	"net"
	"zinx/config"
)

//iServer接口实现，定义一个Server服务类
type Server struct {
	//服务器ip
	IPVersion string
	IP        string
	//服务器port
	Port int
	//服务器名称
	Name string

	//路由属性
	MsgHandler ziface.IMsgHandle


	//链接管理模块
	connMgr ziface.IConnManager

	//该server创建链接之后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)

	//该server销毁链接之前自动调用Hook函数
	OnConnStop func(conn ziface.IConnection)

}

//定义一个具体的回显业务 针对handlefunc
/*func CallBackBusi(request ziface.IRequest) error {
	//回显业务
	fmt.Println("【conn Handle】 CallBack..")
	c := request.GetConnection().GetTCPConnection()
	buf := request.GetData()
	cnt := request.GetDataLen()

	if _,err :=c.Write(buf[:cnt]) ;err!=nil {
		fmt.Println("write Back ",err)
		return err
	}
	return nil

}*/


//初始化的New方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      config.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        config.GlobalObject.Host,
		Port:      config.GlobalObject.Port,
		MsgHandler:NewMsgHandler(),
		connMgr:NewConnManager(),
	}

	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

//启动服务器
//原生socket服务器编程
func (s *Server) Start() {

	fmt.Printf("[start] Sever Linstenner at IP:%s,Port:%d,is starting \n", s.IP, s.Port)

	//启动工作池
	s.MsgHandler.StartWorkerPool()

	//创建套接字：得到一个TCP的addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error", err)
		return
	}
	//监听服务器地址
	listenner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("Listen", s.IPVersion, "err,", err)
	}

	//生成Id的累加器
	var cid uint32
	cid= 0


	//阻塞等待客户端发送请求
	go func() {
		for {
			//永久存在
			//阻塞等待客户端请求
			conn, err := listenner.AcceptTCP()//只是针对TCP协议
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			if s.connMgr.Len() >= config.GlobalObject.MaxConn {
				fmt.Println("---> Too many Connection MAxConn = ", config.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//创建一个Connection对象
			dealConn := NewConnection(s,conn,cid,s.MsgHandler)
			cid++

			go dealConn.Start()


		}
	}()


}






//停止服务器
func (s *Server) Stop() {
	//服务器应清空全部链接
	s.connMgr.CleanConn()

}

//运行服务器
func (s *Server) Serve() {
	//启动server 的监听功能
	s.Start()//并不希望他永久的阻塞

	//做一些其他的扩展
	select {}//main函数不会退出，go routine中for不会结束

}

func (s *Server) AddRouter(msgID uint32,router ziface.IRouter)  {
	s.MsgHandler.AddRouter(msgID ,router )
	fmt.Println("Add Router SUCC!! msgID = ", msgID)
}

func (s *Server)GetConnMgr() ziface.IConnManager {
	return s.connMgr
}


//该server创建链接之后自动调用Hook函数
func (s *Server)AddOnConnStart (hookFunc func(conn ziface.IConnection)){
	s.OnConnStart = hookFunc
}

//该server销毁链接之前自动调用Hook函数
func (s *Server)AddOnConnStop (hookFunc func(conn ziface.IConnection)){
	s.OnConnStop = hookFunc
}

//调用
func (s *Server)CallOnConnStart(conn ziface.IConnection)  {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart()...")
		s.OnConnStart(conn)
	}

}

//调用
func (s *Server)CallOnConnStop(conn ziface.IConnection)  {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop()...")
		s.OnConnStop(conn)
	}

}