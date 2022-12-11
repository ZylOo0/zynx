package znet

import (
	"fmt"
	"net"
)

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	Router      *Router
	ConnMgr     *ConnManager
	OnConnStart func(conn *Connection)
	OnConnStop  func(conn *Connection)
}

func (s *Server) Start() {
	fmt.Printf("[Zynx] Server Name: %s, listening at IP: %s, Port: %d is starting\n",
		config.Name, config.Host, config.TcpPort)
	fmt.Printf("[Zynx] Version %s, MaxConn:%d, MaxDataLen:%d\n",
		config.Version, config.MaxConn, config.MaxDataLen)

	go func() {
		// 开启工作池
		s.Router.StartWorkerPool()

		// TCPAddr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// TCPListener
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Zynx server ", s.Name, " success, Listening...")
		var cid uint32
		cid = 0

		for {
			// TCPConn
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			if s.ConnMgr.Len() >= config.MaxConn {
				fmt.Println("===== Too Many Connections, MaxConn = ", config.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.Router)
			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zynx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddHandler(msgID uint32, handler Handler) {
	s.Router.AddHandler(msgID, handler)
	fmt.Println("Add Handler Success!!")
}

func (s *Server) GetConnMgr() *ConnManager {
	return s.ConnMgr
}

func NewServer() *Server {
	s := &Server{
		Name:      config.Name,
		IPVersion: "tcp4",
		IP:        config.Host,
		Port:      config.TcpPort,
		Router:    NewRouter(),
		ConnMgr:   NewConnManager(),
	}

	return s
}

func (s *Server) SetOnConnStart(hookFunc func(conn *Connection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn *Connection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn *Connection) {
	if s.OnConnStart != nil {
		fmt.Println("----- Call OnConnStart() -----")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn *Connection) {
	if s.OnConnStart != nil {
		fmt.Println("----- Call OnConnStop() -----")
		s.OnConnStop(conn)
	}
}
