package znet

import (
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[uint32]*Connection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]*Connection),
	}
}

func (connMgr *ConnManager) Add(conn *Connection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), " added to ConnManager successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn *Connection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), " removed from ConnManager successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (*Connection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections success! conn num = ", connMgr.Len())
}
