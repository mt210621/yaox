package znet

import (
	"awesomeProject/src/zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

// 连接管理模块

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

// NewConnManager 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源  map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将 conn 加入到connmangerr中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection id = ", conn.GetConnID(), " add to connmanager successfully: conn num =", connMgr.Len())
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源  map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//delete 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connection id = ", conn.GetConnID(), " remove from connmanager successfully: conn num =", connMgr.Len())
}

// Get 根据connid获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源  map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		//找到了
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// Len 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// ClearConn 清除并终止所有的连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源  map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println(" clear all connection succ!!! conn num = ", connMgr.Len())
}