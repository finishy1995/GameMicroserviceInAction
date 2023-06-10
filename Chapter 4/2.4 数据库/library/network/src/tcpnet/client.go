package tcpnet

import (
	"ProjectX/library/log"
	"ProjectX/library/network/codec"
	"ProjectX/library/network/core"
	"ProjectX/library/routine"
	"net"
	"sync"
	"time"
)

// Client 客户端
type Client struct {
	sync.Mutex
	reconnect bool
	addr      string
	isConnect bool
	closeSig  chan bool
	closeFlag bool
	newAgent  core.GetAgent
	codec     core.Codec
	wg        sync.WaitGroup
	conn      *Conn
}

// Start 开启客户端连接
func (client *Client) Start(address string, newAgent core.GetAgent, opts ...core.ClientOption) error {
	if newAgent == nil {
		return core.ErrInvalidGetAgentFunc
	}
	if !core.VerifyAddress(address) {
		return core.ErrInvalidAddress
	}
	client.newAgent = newAgent
	client.addr = address
	options := core.DefaultClientOptions
	for _, o := range opts {
		o(&options)
	}
	client.codec = new(codec.LengthFieldBasedFrameCodec)
	if options.Context != nil {
		if i, ok := options.Context["stick"]; ok {
			if !i.(bool) {
				client.codec = new(codec.BuiltInCodec)
			}
		}
	}

	client.reconnect = options.Reconnect
	client.isConnect = false
	client.closeSig = make(chan bool, 1)
	client.closeFlag = false

	log.Info("TCP trying to connect %s", client.addr)

	return nil
}

// Run 执行主逻辑
func (client *Client) Run() {
	firstStart := make(chan bool, 1)
	firstStart <- true
	for {
		select {
		case <-client.closeSig:
			return
		case <-firstStart:
			break
		case <-time.After(core.DefaultConnectMaxWait):
			break
		}

		conn, err := net.DialTimeout("tcp", client.addr, core.DefaultConnectMaxWait)
		if err != nil {
			continue
		}

		log.Info("TCP connect to %s", client.addr)
		newAgent := client.newAgent()
		if newAgent == nil {
			panic(core.ErrInvalidGetAgentFunc)
		}
		client.Lock()
		tcpConn := pool.Get().(*Conn)
		tcpConn.Init(conn, client.codec)
		client.conn = tcpConn
		tcpConn.setAgent(newAgent)
		client.Unlock()
		client.wg.Add(1)

		runFunc := func() {
			defer func() {
				newAgent.OnClose(tcpConn)
				tcpConn.Close()
				client.isConnect = false
				client.wg.Done()
			}()
			client.isConnect = true
			tcpConn.Run()
		}
		for {
			err = routine.Run(false, runFunc)
			if err != nil {
				time.Sleep(time.Millisecond * 20)
			} else {
				break
			}
		}
		client.wg.Wait()
	}
}

// Close 关闭客户端连接
func (client *Client) Close() {
	if client.closeFlag {
		return
	}
	client.Lock()
	defer client.Unlock()
	if client.closeFlag {
		return
	}
	client.closeSig <- true
	client.closeFlag = true
	if client.conn != nil {
		client.conn.Close()
		client.wg.Wait()
	}
}

// IsConnected 是否处于连接状态
func (client *Client) IsConnected() bool {
	return client.isConnect
}
