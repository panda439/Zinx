package znet

import (
	"zinx/ziface"
	"bytes"
	"encoding/binary"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}


//获取二进制包的头部长度 固定返回8
func (dp *DataPack)GetHeadLen() uint32{
	return 8
}


//封包方法 -- 将Message打包成 |datalen|dataID|data|

func (dp *DataPack)Pack(msg ziface.IMessage)([]byte, error){
	//创建一个存放二进制的字节缓冲
	dataBuffer := bytes.NewBuffer([]byte{})

	//将datalen写进buffer中 小头传递
	if err:=binary.Write(dataBuffer,binary.LittleEndian,msg.GetMsgLen());err!=nil{
		return nil,err
	}
	//将id写进buffer中
	if err:=binary.Write(dataBuffer,binary.LittleEndian,msg.GetMsgId());err!=nil{
		return nil,err
	}
	//将 data写进buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgData()) ; err != nil {
		return nil, err
	}
	//返回这个缓冲

	return dataBuffer.Bytes(),nil
}
//拆包方法 --  |datalen|dataID|data| 拆解到Message结构体中
func (dp *DataPack)UnPack(binaryData []byte)(ziface.IMessage,error){
	//解包的时候 分两次解压，第一次读取固定的长度8字节，第二次是根据len再次进行read
	//uint32 4字节读取满后再往下进行
	msgHead := &Message{}
	
	//创建一个读取二进制数据流的io.Reader
	dataBuff := bytes.NewReader(binaryData)
	
	//将二进制流 先读datalen放在msg的DataLen属性中
	if err := binary.Read(dataBuff, binary.LittleEndian, &msgHead.Len); err!=nil {
		return nil, err
	}
	

	//将二进制流的DataID 放在Msg的DataID属性中
	if err := binary.Read(dataBuff, binary.LittleEndian, &msgHead.Id); err != nil {
		return nil, err
	}


	return msgHead,nil

}