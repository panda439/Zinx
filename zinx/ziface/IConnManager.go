package ziface

type IConnManager interface {
	Add(conn IConnection)

	Remove(connID uint32)

	Get(connID uint32) (IConnection, error)

	Len() uint32

	CleanConn()


}
