package tcpgnet

import (
	"ProjectX/library/log"
	"ProjectX/library/network/core"
	"encoding/binary"
	"sync"
	"time"

	"github.com/panjf2000/gnet"
)

// Server TCP GNet 服务器
type Server struct {
	*gnet.EventServer
	// 连接管理
	connSet map[string]*Conn

	addr       string
	maxConnNum int
	newAgent   core.GetAgent
	codec      gnet.ICodec

	connMutex sync.RWMutex
	wgConn    sync.WaitGroup
	closeFlag bool
}

func defaultCodec() gnet.ICodec {
	ec := gnet.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}
	dc := gnet.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		LengthFieldOffset:   0,
		InitialBytesToStrip: 4,
	}
	return gnet.NewLengthFieldBasedFrameCodec(ec, dc)
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
	server.codec = defaultCodec()
	if options.Context != nil {
		if i, ok := options.Context["stick"]; ok {
			if !i.(bool) {
				server.codec = new(gnet.BuiltInFrameCodec)
			}
		}
	}

	// 初始化数组
	server.connSet = make(map[string]*Conn)
	server.closeFlag = false

	log.Info("TCP Listen %s", server.addr)

	return nil
}

// Run 执行服务端逻辑
func (server *Server) Run() {
	// 创建监听
	err := gnet.Serve(server, "tcp://"+server.addr,
		gnet.WithMulticore(true),
		gnet.WithLogger(new(logger)),
		gnet.WithReusePort(false),
		gnet.WithTicker(true),
		gnet.WithTCPKeepAlive(core.DefaultKeepAlive),
		gnet.WithCodec(server.codec),
		gnet.WithLockOSThread(true),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
	)
	if err != nil {
		panic(err)
	}
}

// Close 关闭TCP监听
func (server *Server) Close() {
	if server.closeFlag {
		return
	}
	server.connMutex.RLock()
	if server.closeFlag {
		return
	}
	server.closeFlag = true
	for _, conn := range server.connSet {
		conn.Close()
	}
	server.connMutex.RUnlock()
	log.Info("TCP Close %s", server.addr)
	server.wgConn.Wait()
}

// GetConnNum 获取所有连接的数量
func (server *Server) GetConnNum() (num int) {
	return len(server.connSet)
}

// OnOpened 当有新连接建立时调用
func (server *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	agent := server.newAgent()
	if agent == nil {
		action = gnet.Close
		log.Error("New agent error: %v", core.ErrInvalidGetAgentFunc)
		return
	}
	tcpConn := pool.Get().(*Conn)
	tcpConn.Init(c)

	server.connMutex.Lock()
	// 如果超过了最大限制，则关闭连接
	if server.GetConnNum() >= server.maxConnNum {
		server.connMutex.Unlock()
		pool.Put(tcpConn)
		action = gnet.Close
		log.Info("Over connection limit!")
		return
	}
	server.connSet[c.RemoteAddr().String()] = tcpConn
	server.connMutex.Unlock()

	tcpConn.setAgent(agent)
	server.wgConn.Add(1)
	return
}

// OnClosed 当连接关闭时调用
func (server *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	server.connMutex.RLock()
	if conn, ok := server.connSet[c.RemoteAddr().String()]; ok {
		server.connMutex.RUnlock()
		if !server.closeFlag {
			// 如果是服务器还没关闭的情况下关闭了链接
			conn.Close()
			if err != nil {
				log.Info("Connection closed err %v", err)
			}
			server.connMutex.Lock()
			delete(server.connSet, c.RemoteAddr().String())
			server.connMutex.Unlock()
		}
		pool.Put(conn)
		server.wgConn.Done()
	} else {
		server.connMutex.RUnlock()
		// 当连接数超过上限时会触发
		log.Error("Connection closed when not store in map")
	}

	return
}

// React 当有消息收到时调用
func (server *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	server.connMutex.RLock()
	if conn, ok := server.connSet[c.RemoteAddr().String()]; ok {
		server.connMutex.RUnlock()
		conn.onMessage(frame)
	} else {
		server.connMutex.RUnlock()
	}
	return
}

// Tick 每隔一段时间调用
func (server *Server) Tick() (delay time.Duration, action gnet.Action) {
	if server.closeFlag {
		action = gnet.Shutdown
		return
	}
	delay = core.UpdateInterval

	return
}
