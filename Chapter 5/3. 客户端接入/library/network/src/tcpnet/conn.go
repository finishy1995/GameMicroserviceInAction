package tcpnet

import (
	"ProjectX/library/log"
	"ProjectX/library/network/codec"
	"ProjectX/library/network/core"
	"ProjectX/library/network/protocol"
	"bytes"
	"net"
	"sync"
	"time"
)

var (
	heartBeatMsg = []byte("|")
	protoc       = protocol.ProtocolV001
)

// Conn tcp 连接
type Conn struct {
	sync.Mutex
	codec.ConnHelper

	id        core.ID
	conn      net.Conn
	closeFlag bool
	closeSig  chan bool
	agent     core.Agent
	codec     core.Codec

	// 最近一次心跳包时间
	lastHeartbeatTime int64
	// 最近一次收到包时间
	lastRecvTime int64
}

func getTime() int64 {
	return time.Now().UnixNano() / 1000000
}

// Init 初始化
func (tcpConn *Conn) Init(conn net.Conn, codec core.Codec) {
	if codec == nil {
		panic(core.ErrInvalidCodec)
	}
	tcpConn.conn = conn
	tcpConn.InitBuffer()
	tcpConn.codec = codec
	t := getTime()
	tcpConn.lastHeartbeatTime = t
	tcpConn.lastRecvTime = t
	tcpConn.closeSig = make(chan bool, 1)
	tcpConn.id = core.GenerateID()
}

// Close 断连
func (tcpConn *Conn) Close() {
	if tcpConn.closeFlag {
		return
	}
	tcpConn.Lock()
	defer tcpConn.Unlock()

	if tcpConn.closeFlag {
		return
	}
	tcpConn.closeSig <- true
	tcpConn.closeFlag = true
	err := tcpConn.conn.(*net.TCPConn).SetLinger(0)
	if err != nil {
		log.Error("tcp close failed, error: %s", err.Error())
	}
	err = tcpConn.conn.Close()
	if err != nil {
		panic(err)
	}

	tcpConn.conn = nil
	tcpConn.agent = nil
	tcpConn.codec = nil
}

// Write b 必须在其他协程中不被修改
func (tcpConn *Conn) Write(b []byte) (n int, err error) {
	if tcpConn.closeFlag || b == nil {
		return
	}
	tcpConn.Lock()
	defer tcpConn.Unlock()
	if tcpConn.closeFlag {
		return
	}

	out, _ := tcpConn.codec.Encode(tcpConn, b)
	return tcpConn.conn.Write(out)
}

// LocalAddr 本地socket端口地址
func (tcpConn *Conn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

// RemoteAddr 远程socket端口地址
func (tcpConn *Conn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}

// setAgent 设置 Agent
func (tcpConn *Conn) setAgent(agent core.Agent) {
	tcpConn.agent = agent
	tcpConn.agent.OnConnect(tcpConn)
}

// Run 循环运行
func (tcpConn *Conn) Run() {
	b := make([]byte, core.TCPMaxPackageSize)
	t := time.Millisecond * 2
	update := core.UpdateInterval - t
	for {
		if tcpConn.closeFlag {
			return
		}
		err := tcpConn.conn.SetReadDeadline(time.Now().Add(t))
		if err != nil {
			return
		}
		n, err := tcpConn.conn.Read(b)
		now := getTime()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				//if now - tcpConn.lastRecvTime > core.TimeoutTime {
				//	return
				//}
				//if now - tcpConn.lastHeartbeatTime > core.HeartbeatTime {
				//	_, err := tcpConn.Write(heartBeatMsg)
				//	tcpConn.lastHeartbeatTime = now
				//	if err != nil {
				//		return
				//	}
				//}
				continue
			} else if ok && netErr.Temporary() {
				select {
				case <-tcpConn.closeSig:
					return
				case <-time.After(update):
					continue
				}
			} else {
				return
			}
		}

		if n > 0 {
			if tcpConn.agent == nil {
				panic(core.ErrInvalidAgent)
			}
			tcpConn.lastRecvTime = now
			tcpConn.lastHeartbeatTime = now
			tcpConn.PushPacket(b[:n])

			for {
				out, err := tcpConn.codec.Decode(tcpConn)
				if err != nil {
					if err == core.ErrPacketSplit {
						break
					} else {
						panic(err)
					}
				}
				if out == nil {
					continue
				}
				//if len(out) == len(heartBeatMsg) && out[0] == heartBeatMsg[0] {
				//	// 这个是心跳包，应用层不处理
				//	continue
				//}
				if bytes.Equal(out, protoc.ReceiveMsg) {
					tcpConn.Write(protoc.ReplyMsg)
					continue
				}

				tcpConn.agent.OnMessage(out, tcpConn)
			}
			continue
		} else {
			select {
			case <-tcpConn.closeSig:
				return
			case <-time.After(update):
				continue
			}
		}
	}
}
