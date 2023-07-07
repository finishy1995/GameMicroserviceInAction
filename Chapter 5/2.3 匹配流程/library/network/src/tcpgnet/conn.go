package tcpgnet

import (
	"ProjectX/library/network/core"
	"ProjectX/library/network/protocol"
	"bytes"
	"net"
	"sync"

	"github.com/panjf2000/gnet"
)

var (
	protoc = protocol.ProtocolV001
)

type Conn struct {
	sync.Mutex
	gnetConn  gnet.Conn
	closeFlag bool
	agent     core.Agent
}

// Init 初始化
func (conn *Conn) Init(c gnet.Conn) {
	conn.gnetConn = c
	conn.closeFlag = false
}

// Run 主逻辑
func (conn *Conn) Run() {
}

// Close 关闭
func (conn *Conn) Close() {
	if conn.closeFlag {
		return
	}
	conn.Lock()
	defer conn.Unlock()
	if conn.closeFlag {
		return
	}
	conn.closeFlag = true
	if conn.agent != nil {
		conn.agent.OnClose(conn)
	}
	if conn.gnetConn != nil {
		err := conn.gnetConn.Close()
		if err != nil {
			panic(err)
		}
	}

	conn.agent = nil
	conn.gnetConn = nil
}

// Write 发送数据
func (conn *Conn) Write(b []byte) (n int, err error) {
	if conn.closeFlag {
		return
	}
	conn.Lock()
	defer conn.Unlock()
	if conn.closeFlag {
		return
	}
	n = len(b)
	err = conn.gnetConn.AsyncWrite(b)
	return
}

// LocalAddr 本地地址
func (conn *Conn) LocalAddr() net.Addr {
	return conn.gnetConn.LocalAddr()
}

// RemoteAddr 远程地址
func (conn *Conn) RemoteAddr() net.Addr {
	return conn.gnetConn.RemoteAddr()
}

func (conn *Conn) setAgent(agent core.Agent) {
	conn.agent = agent
	conn.agent.OnConnect(conn)
}

func (conn *Conn) onMessage(b []byte) {
	if conn.agent != nil && !conn.closeFlag {
		if bytes.Equal(b, protoc.ReceiveMsg) {
			conn.Write(protoc.ReplyMsg)
			return
		}

		conn.agent.OnMessage(b, conn)
	}
}
