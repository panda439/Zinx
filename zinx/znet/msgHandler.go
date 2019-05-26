package znet

import (
	"zinx/ziface"
	"fmt"
)

type MsgHandler struct {
	//存放路由集合的map
	Apis map[uint32] ziface.IRouter //就是开发者全部的业务，消息ID和业务

}

func NewMsgHandler() ziface.IMsgHandle {
	return &MsgHandler{
		Apis:make(map[uint32]ziface.IRouter),
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
}