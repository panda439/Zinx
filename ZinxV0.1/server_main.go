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

type HelloRouter struct {
	znet.BaseRouter
}

func (this *PingRouter)Handle(request ziface.IRequest)  {
	fmt.Println("Call Router Handler...")
	//给客户端回写一个 数据
	err := request.GetConnection().Send(1,[]byte("200ping...ping...ping"))
	if err!=nil{
		fmt.Println(err)
	}
}
func (this *HelloRouter)Handle(request ziface.IRequest)  {
	fmt.Println("Call Router Handler...")
	//给客户端回写一个 数据
	err := request.GetConnection().Send(2,[]byte("201ping...ping...ping"))
	if err!=nil{
		fmt.Println(err)
	}
}




func main() {
	//创建一个zinx server对象
	s := znet.NewServer("zinx v0.1")

	//注册一些自定义的业务
	s.AddRouter(1,&PingRouter{} )
	s.AddRouter(2,&HelloRouter{} )

	//让server对象 启动服务
	s.Serve()
}
