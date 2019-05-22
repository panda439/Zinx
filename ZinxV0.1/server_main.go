/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/22 12:16
* @Mail: danbing.at@gmail.com
*/
package main

import "zinx/znet"

func main() {
	//创建一个zinx server对象
	s := znet.NewServer("zinx v0.1")

	//注册一些自定义的业务


	//让server对象 启动服务
	s.Serve()
}
