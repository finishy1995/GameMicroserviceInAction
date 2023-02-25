package tcpnet

import (
	"ProjectX/library/log"
	"ProjectX/library/network/codec"
	"ProjectX/library/network/core"
	"ProjectX/library/routine"
	"crypto/tls"
	"net"
	"sync"
	"time"
)

// Server tcp 服务器
type Server struct {
	// 连接管理
	connSet map[core.ID]*Conn

	addr       string
	tls        bool
	certFile   string
	keyFile    string
	maxConnNum int
	newAgent   core.GetAgent
	codec      core.Codec

	ln        net.Listener
	connMutex sync.Mutex
	wgLn      sync.WaitGroup
	wgConn    sync.WaitGroup
	closeSig  chan bool
	closeFlag bool
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
	server.codec = new(codec.LengthFieldBasedFrameCodec)
	// 定制化设置
	if options.Context != nil {
		if t, ok := options.Context["tls"]; ok {
			server.tls = t.(bool)
		}
		if certFile, ok := options.Context["certFile"]; ok {
			server.certFile = certFile.(string)
		}
		if keyFile, ok := options.Context["keyFile"]; ok {
			server.keyFile = keyFile.(string)
		}
		if i, ok := options.Context["stick"]; ok {
			if !i.(bool) {
				server.codec = new(codec.BuiltInCodec)
			}
		}
	}

	// 初始化数组
	server.connSet = make(map[core.ID]*Conn)
	server.closeSig = make(chan bool, 1)
	server.closeFlag = false

	// 创建监听
	ln, err := net.Listen("tcp", server.addr)
	if err != nil {
		return err
	}

	// 创建 tls
	if server.tls {
		tlsConf := new(tls.Config)
		tlsConf.Certificates = make([]tls.Certificate, 1)
		tlsConf.Certificates[0], err = tls.LoadX509KeyPair(server.certFile, server.keyFile)
		if err == nil {
			ln = tls.NewListener(ln, tlsConf)
			log.Info("TCP Listen TLS load success")
		} else {
			return err
		}
	}

	server.ln = ln
	log.Info("TCP Listen %s", server.addr)

	return nil
}

// Run 执行服务端逻辑
func (server *Server) Run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()

	var tempDelay time.Duration
	for {
		if len(server.closeSig) > 0 {
			return
		}

		if server.GetConnNum() >= server.maxConnNum {
			select {
			case <-server.closeSig:
				return
			case <-time.After(core.DefaultConnectWaitStart):
				continue
			}
		}

		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = core.DefaultConnectWaitStart
				} else {
					tempDelay *= core.DefaultConnectWaitMut
				}
				if tempDelay > core.DefaultConnectMaxWait {
					tempDelay = core.DefaultConnectMaxWait
				}
				log.Error("accept error: %v; retrying in %v", err, tempDelay)
				select {
				case <-server.closeSig:
					return
				case <-time.After(tempDelay):
					continue
				}
			} else {
				return
			}
		}
		tempDelay = 0

		agent := server.newAgent()
		if agent == nil {
			panic(core.ErrInvalidGetAgentFunc)
		}
		tcpConn := pool.Get().(*Conn)
		tcpConn.Init(conn, server.codec)
		tcpConnID := tcpConn.id
		server.connSet[tcpConnID] = tcpConn
		tcpConn.setAgent(agent)
		server.wgConn.Add(1)

		closeConn := func() {
			agent.OnClose(tcpConn)
			tcpConn.Close()

			if !server.closeFlag {
				server.connMutex.Lock()
				delete(server.connSet, tcpConnID)
				server.connMutex.Unlock()
			}
			server.wgConn.Done()
		}

		err = routine.Run(true, func() {
			defer closeConn()
			tcpConn.Run()
		})
		if err != nil {
			closeConn()
		}
	}
}

// Close 关闭TCP监听
func (server *Server) Close() {
	if server.closeFlag {
		return
	}
	server.connMutex.Lock()
	server.closeFlag = true
	server.closeSig <- true
	for _, conn := range server.connSet {
		conn.Close()
	}
	server.connMutex.Unlock()
	err := server.ln.Close()
	log.Info("TCP Close %s", server.addr)
	if err != nil {
		panic(err)
	}
	server.wgLn.Wait()
	server.wgConn.Wait()
}

// GetConnNum 获取所有连接的数量
func (server *Server) GetConnNum() (num int) {
	return len(server.connSet)
}
