package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	TcpServer    *Server
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	ExitChan     chan bool   // 用于告知当前连接已经停止(由 Reader 告知 Writer)
	msgChan      chan []byte // 用于读写协程的通信
	Router       *Router
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server *Server, conn *net.TCPConn, connID uint32, router *Router) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		Router:    router,
		isClosed:  false,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property:  make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c)

	return c
}

var reqID uint32 = 0

func (c *Connection) StartReader() {
	fmt.Println("[Reader GoRoutine is running...]")
	defer fmt.Println("[Reader exits!], connID = ", c.ConnID, ", remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		msg, err := c.ReadMsg()
		if err != nil {
			fmt.Println("Connection ", c.ConnID, " receive msg error")
			return
		}

		req := Request{
			id:   reqID,
			conn: c,
			msg:  msg,
		}
		reqID++

		if config.WorkerPoolSize > 0 {
			c.Router.SendRequestToTaskQueue(&req)
		} else {
			fmt.Println("worker pool size < 0")
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer GoRoutine is running...]")
	defer fmt.Println("[conn Writer exits!]", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	go c.StartReader()
	go c.StartWriter()

	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	err := c.Conn.Close()
	if err != nil {
		fmt.Println("TCP Connection Close() error, ", err)
	}

	// 告知 Writer 关闭
	c.ExitChan <- true

	c.TcpServer.GetConnMgr().Remove(c)

	close(c.msgChan)
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsgToWriter(msgId uint32, data []byte) error { // 由用户使用
	if c.isClosed == true {
		return errors.New("connection is closed when sending msg")
	}

	msg := NewMessage(msgId, data)
	binaryMsg, err := Pack(msg)
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
	}

	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) ReadMsg() (*Message, error) {
	head := make([]byte, HeadLength)
	if _, err := io.ReadFull(c.GetTCPConnection(), head); err != nil {
		fmt.Println("read msg head error", err)
		return nil, err
	}

	msg, err := Unpack(head)
	if err != nil {
		fmt.Println("unpack error", err)
	}

	var data []byte
	if msg.GetDataLen() > 0 {
		data = make([]byte, msg.GetDataLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
			fmt.Println("read msg data error", err)
		}
	}
	msg.SetData(data)

	return msg, nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
