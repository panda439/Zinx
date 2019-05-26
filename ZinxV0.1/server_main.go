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

//创建链接之后的执行的钩子函数
func DoConntionBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConntionBegin  ....")
	//链接一旦创建成功 给用户返回一个消息
	if err := conn.Send(202, []byte("Hello welcome to zinx...")); err !=nil {
		fmt.Println(err)
	}
}

//链接销毁之前执行的钩子函数
func DoConntionLost(conn ziface.IConnection) {
	fmt.Println("===> DoConntionLost  ....")
	fmt.Println("Conn id ", conn.GetConnID(), "is Lost!.")
}




func main() {
	//创建一个zinx server对象
	s := znet.NewServer("zinx v0.1")

	//注册一个创建链接之后的方法业务
	s.AddOnConnStart(DoConntionBegin)
	//注册一个链接断开之前的方法业务
	s.AddOnConnStop(DoConntionLost)

	//注册一些自定义的业务
	s.AddRouter(1,&PingRouter{} )
	s.AddRouter(2,&HelloRouter{} )

	//让server对象 启动服务
	s.Serve()
}
