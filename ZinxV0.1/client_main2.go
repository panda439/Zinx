package main

import (
	"fmt"
	"time"
	"net"
	"zinx/znet"
	"io"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)

	//之前connect服务器得到一个已经建立好的conn句柄
	conn,err := net.Dial("tcp",":8999")
	if err != nil {
		fmt.Println("client start err",err)
		return
	}
	for  {
		dp := znet.NewDataPack()

		binaryMsg,err := dp.Pack(znet.NewMsgPackage(2,[]byte("Zinx 0.6 client Test Massage")))

		if err!=nil {
			fmt.Println("Pack error",err)
			return
		}
		if _,err = conn.Write(binaryMsg);err!=nil {
			fmt.Println(err)
			return
		}
		binareyHead := make([]byte,dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binareyHead); err != nil {
			fmt.Println(err)
			return
		}

		msgHead,err := dp.UnPack(binareyHead)
		if msgHead.GetMsgLen()>0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte,msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error",err)
				return
			}

			fmt.Println("---> Recv Server Msg : id =",msg.Id,"len:",msg.Len,"data = ",string(msg.Data))

		}



		time.Sleep(1 *time.Second)



	}
}
