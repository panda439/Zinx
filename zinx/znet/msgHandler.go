package znet

import (
	"zinx/ziface"
	"fmt"
	"zinx/config"
)

type MsgHandler struct {
	//存放路由集合的map
	Apis map[uint32] ziface.IRouter //就是开发者全部的业务，消息ID和业务

	//负责Worker取任务的消息队列 一个worker对应一个任务队列
	TaskQueue []chan ziface.IRequest

	//worker工作池的worker的数量
	WorkerPoolSize uint32


}

func NewMsgHandler() ziface.IMsgHandle {
	return &MsgHandler{
		Apis:make(map[uint32]ziface.IRouter),
		TaskQueue:make([]chan ziface.IRequest,config.GlobalObject.WorkerPoolSize),
		WorkerPoolSize:config.GlobalObject.WorkerPoolSize,
	}
}


func (mh *MsgHandler)AddRouter(msgID uint32,router ziface.IRouter){
	//判断新添加的msgID是否存在
	if _ ,ok := mh.Apis[msgID];ok{
		//msgId已经注册
		fmt.Println("repeat Api msgID = ",msgID)
		return
	}
	//添加msgID和 router的对应关系
	mh.Apis[msgID] = router
	fmt.Println("Apd api MsgID = ",msgID,"succ!")

}

//调度路由，根据MsgID
func (mh *MsgHandler)DoMsgHandler(request ziface.IRequest) {
	//从Request取到MsgiD
	router, ok := mh.Apis[request.GetMsg().GetMsgId()]
	if !ok {
		fmt.Println("api MsgID=", request.GetMsg().GetMsgId(), "Not Found! Need Add!")
	}
	//根据msgID找到对应的router 进行调用
	router.Handle(request)
	router.PreHandle(request)
	router.PostHandle(request)
}

func (mh *MsgHandler) startOneWorker(workerID int,taskQueue chan ziface.IRequest) {
	fmt.Println("workerID = ",workerID,"is starting...")
	//不断的从对应的管道等待数据
	for  {
		select {
		case req:= <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}


//启动Worker工作池
func (mh *MsgHandler)StartWorkerPool(){
	fmt.Println("WorkPool is started..")

	//根据WorkerPoolSize创建worker groutin
	for i := 0;i<int(mh.WorkerPoolSize); i++ {


		mh.TaskQueue[i] = make(chan ziface.IRequest,config.GlobalObject.MaxWorkerTaskLen)
		//启动一个Worker ，阻塞等待消息从对应的管道中进来

		go mh.startOneWorker(i,mh.TaskQueue[i])

	}


}
//将消息添加到worker工作池中
//Reader来调用
func (mh *MsgHandler)SendMsgToTaskQueue(request ziface.IRequest){

	//将消息平均分配给worker确定当前的request到底给哪个worker来处理
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//直接将request发送给对应的worker的taskqueue
	mh.TaskQueue[workerID] <- request



}