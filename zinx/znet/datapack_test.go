package znet

import (
	"testing"
	"fmt"
	"net"
	"io"
)


//函数名 Test 开头 后面的函数名 自定义
//行参（t *testing.T）
func TestDataPack(t *testing.T)  {
	fmt.Println("nihao I am testing...")

	listenner,err:= net.Listen("tcp",":7777")
	if err != nil {
		fmt.Println("server listenner err",err)
		return
	}

	go func() {

		conn,err := listenner.Accept()
		if err != nil {
			fmt.Println("server accpet err",err)
		}


		//读写业务
		go func(conn *net.Conn) {
			//读取客户端请求
			//---拆包过程---
			//|datalen|id|data|
			dp := NewDataPack()
			for  {
				headData := make([]byte,dp.GetHeadLen())//8字节
				_,err := io.ReadFull(*conn,headData)//直到headData填充满才会返回，否则阻塞
				if err != nil {
					fmt.Println("read head error")
					break
				}
				msgHead,err := dp.UnPack(headData)
				if err != nil {
					fmt.Println("",err)
					return
				}
				if msgHead.GetMsgLen()>0 {
					//有内容，需要进行二次读取
					//将msgHead 进行向下转换 将iMessage转换成Message
					msg:= msgHead.(*Message)
					//开辟
					msg.Data = make([]byte,msg.GetMsgLen())
					//根据长度读取
					_,err:=io.ReadFull(*conn,msg.Data)
					if err != nil {
						fmt.Println("server unpack data error",err)
						return
					}
					fmt.Println("recv MsgId = ",msg.Id,"datalen = ",msg.Len,"data = ",string(msg.Data))
				}

			}

		}(&conn)
	}()
	conn,err := net.Dial("tcp",":7777")
	if err != nil {
		fmt.Println("clien dail err:",err)
		return
	}
	dp := NewDataPack()
	msg1 := &Message{
		Id:1,
		Len:4,
		Data:[]byte{'z','i','n','x'},
	}
	sendData1,err := dp.Pack(msg1)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg2 := &Message{
		Id:1,
		Len:5,
		Data:[]byte{'h','e','l','l','o'},
	}
	sendData2,err := dp.Pack(msg2)
	if err != nil {
		fmt.Println(err)
		return
	}

	//将两个包粘在一起
	sendData1  = append(sendData1,sendData2...)//打散？？？

	//发送
	conn.Write(sendData1)

	//让test不结束
	select {}
}
