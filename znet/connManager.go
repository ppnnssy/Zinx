package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinxProject/ziface"
)

/*
链接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex //保护链接集合的读写锁

}

//初始化方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

/*
各种方法的实现
*/

//添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源需要加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn //根据conn的ID存储到map中
	fmt.Println("connID:", conn.GetConnID(), "add to ConnManager successful:conn num:", conn.GetConnID())

}

//删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源需要加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID:", conn.GetConnID(),
		"remove from ConnManager successful:conn num:", conn.GetConnID())

}

//根据链接ID获取链接
func (connMgr *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	//保护共享资源需要加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	conn, ok := connMgr.connections[connId]
	if ok {
		return conn, nil
	} else {
		return nil, errors.New("connection is not FOUND!")
	}
}

//得到当前的链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections) //返回整个map的大小
}

//清楚并终止所有链接。用于回收资源
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源需要加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止链接
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections success!conn num:", connMgr.Len())
}
