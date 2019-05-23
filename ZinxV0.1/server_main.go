/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/22 12:16
* @Mail: danbing.at@gmail.com
*/
package main

import (
	"zinx/znet"
	"zinx/ziface"
	"fmt"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter)PreHandle(request ziface.IRequest)  {
	fmt.Println("Call Router PreHandler...")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}
func (this *PingRouter)Handle(request ziface.IRequest)  {
	fmt.Println("Call Router Handler...")
	//给客户端回写一个 数据
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping..ping..ping...\n"))
	if err != nil {
		fmt.Println("call  ping error")
	}
}
func (this *PingRouter)PostHandle(request ziface.IRequest)  {
	fmt.Println("Call Router PostHandler...")
	//给客户端回写一个 数据
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}


func main() {
	//创建一个zinx server对象
	s := znet.NewServer("zinx v0.1")

	//注册一些自定义的业务
	s.AddRouter(&PingRouter{} )

	//让server对象 启动服务
	s.Serve()
}
