package znet

import (
	"zinx/ziface"
	"sync"
	"fmt"
	"errors"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() ziface.IConnManager  {
	return &ConnManager{
		connections:make(map[uint32]ziface.IConnection),
	}
}


func (connMgr *ConnManager)Add(conn ziface.IConnection){
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("Add connid = ", conn.GetConnID(), "to manager succ!!")

}

func (connMgr *ConnManager)Remove(connID uint32){
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections,connID)
	fmt.Println("Remove connid = ", connID, " from manager succ!!")

}

func (connMgr *ConnManager)Get(connID uint32) (ziface.IConnection, error){
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn,ok := connMgr.connections[connID];ok {
		//找到了
		return conn,nil
	}else {
		return nil,errors.New("connection not FOUND!")
	}


}

func (connMgr *ConnManager)Len() uint32{
	return uint32(len(connMgr.connections))
}

func (connMgr *ConnManager)CleanConn(){
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID,conn := range connMgr.connections {
		conn.Stop()
		//删除链接
		delete(connMgr.connections,connID)
	}
	fmt.Println("Clear All Conections succ! conn num = ", connMgr.Len())
}