package ziface

//定义服务器接口
type IServer interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//开启业务服务方法
	Serve()

	//添加路由的方法
	AddRouter(msgID uint32,router IRouter)

	//提供一个得到连接管理模块的方法
	GetConnMgr() IConnManager

	AddOnConnStart(hookFunc func(conn IConnection))

	AddOnConnStop(hookFunc func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)

}
