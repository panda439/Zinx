package ziface

type DataPack interface {
	GetHeadLen() uint32


	//封包方法 -- 将Message打包成 |datalen|dataID|data|

	Pack(msg IMessage)([]byte, error)
	//拆包方法 --  |datalen|dataID|data| 拆解到Message结构体中
	UnPack()(IMessage,error)

}
