package websocket

import (
	"ProjectX/library/log"
	"ProjectX/library/network/core"
	"ProjectX/library/routine"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	s *http.Server

	addr         string
	newAgent     core.GetAgent
	maxConnNum   int
	connMap      map[*websocket.Conn]*Conn
	connMapMutex sync.Mutex
}

// Start 开始tcp监听
func (server *Server) Start(address string, newAgent core.GetAgent, opts ...core.ServerOption) error {
	// 读取并初始化参数
	if newAgent == nil {
		return core.ErrInvalidGetAgentFunc
	}
	if !core.VerifyAddress(address) {
		return core.ErrInvalidAddress
	}
	server.newAgent = newAgent
	server.addr = address
	options := core.DefaultServerOptions
	for _, o := range opts {
		o(&options)
	}
	if options.MaxConnNum < 0 {
		server.maxConnNum = core.DefaultMaxConnNum
	} else {
		server.maxConnNum = options.MaxConnNum
	}
	server.connMap = make(map[*websocket.Conn]*Conn, 0)
	log.Info("WebSocket Listen %s", server.addr)

	return nil
}

func (server *Server) onWebsocket(w http.ResponseWriter, r *http.Request) {
	u := websocket.NewUpgrader()
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	u.OnOpen(server.OnOpen)
	u.OnMessage(server.OnMessage)
	u.OnClose(server.OnClose)
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		log.Error("SetReadDeadline:", err)
		return
	}
}

// Run 执行服务端逻辑
func (server *Server) Run() {
	// 创建监听
	mux := &http.ServeMux{}
	mux.HandleFunc("/", server.onWebsocket)
	server.s = &http.Server{
		Addr:    server.addr,
		Handler: mux,
	}
	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		panic(err)
	}

	listener = netutil.LimitListener(listener, server.maxConnNum)

	err = routine.Run(false, func() {
		defer listener.Close()

		err := server.s.Serve(listener)
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}
}

// GetConnNum 获取所有连接的数量
func (server *Server) GetConnNum() int {
	return len(server.connMap)
}

func (server *Server) OnMessage(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
	if c, ok := server.connMap[conn]; ok {
		c.agent.OnMessage(data, c)
		return
	}
}

func (server *Server) OnOpen(conn *websocket.Conn) {
	agent := server.newAgent()
	if agent == nil {
		err := conn.Close()
		if err != nil {
			log.Error("agent is nil, close conn error:", err)
		}
		return
	}
	c := newConn(conn, agent)
	server.connMapMutex.Lock()
	server.connMap[conn] = c
	server.connMapMutex.Unlock()
	agent.OnConnect(c)
}

func (server *Server) OnClose(conn *websocket.Conn, err error) {
	c := server.connMap[conn]
	if c == nil {
		return
	}

	c.Close()
	server.connMapMutex.Lock()
	delete(server.connMap, conn)
	server.connMapMutex.Unlock()
}

func (server *Server) Close() {
	if server.s == nil {
		return
	}

	err := server.s.Close()
	if err != nil {
		log.Error("close server error:", err)
	}
	server.s = nil
}
