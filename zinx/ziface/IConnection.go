package ziface

import "net"

type IConnection interface {
	//启动链接
	Start()
	//停止链接
	Stop()
	//获取链接ID
	GetConnID() uint32
	//获取conn的原生socket套接字
	GetTCPConnection() *net.TCPConn

	//获取远程客户端的ip地址
	GetRemoteAddr() net.Addr

	//发送数据给对方客户端
	Send(msgId uint32,msgData []byte)  error

}

//业务处理方法 抽象定义

type HandleFunc func (request IRequest)error