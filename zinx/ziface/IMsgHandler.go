package ziface


//抽象的消息管理模块
type IMsgHandle interface {
	//添加路由
	AddRouter(msgID uint32,router IRouter)

	//调度路由，根据MsgID
	DoMsgHandler(request IRequest)


	//启动Worker工作池
	StartWorkerPool()
	//将消息添加到worker工作池中
	SendMsgToTaskQueue(request IRequest)


}

