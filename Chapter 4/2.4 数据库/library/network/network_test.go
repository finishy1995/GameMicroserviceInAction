package network

import (
	"ProjectX/library/log"
	"ProjectX/library/network/agent"
	"ProjectX/library/network/core"
	"ProjectX/library/routine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

const (
	TestThread      = 10
	TestWaitThread  = 2
	SendMsgNum      = 10
	SendInterval    = time.Millisecond * 10
	WaitMsgSendTime = SendInterval * SendMsgNum * 5
	TestPort1       = "6688"

	DestroyAllowWaitTime = time.Millisecond * 20
	ListenAllowWaitTime  = time.Millisecond * 10
	WaitConnectTime      = time.Millisecond * 20
)

var (
	count      int
	agt        = new(testSendAgent)
	countMutex = new(sync.Mutex)
)

type testSendAgent struct {
	agent.SingleAgent
}

func newTestSendAgent() core.Agent {
	return agt
}

func (agent *testSendAgent) OnConnect(conn core.Conn) {
	routine.Run(false, func() {
		for i := 0; i < SendMsgNum; i++ {
			time.Sleep(SendInterval)
			_, err := conn.Write([]byte("hello"))
			if err != nil {
				log.Error(err.Error())
			}
		}
	})
}

func (agent *testSendAgent) OnMessage(b []byte, conn core.Conn) {
	countMutex.Lock()
	defer countMutex.Unlock()
	count++
	if count == TestThread*10 {
		log.Info("finish message recv")
	}
}

func destroyAfterTest() {
	DestroyAll()
	time.Sleep(DestroyAllowWaitTime)
}

// 测试监听同一个端口
func TestListen(t *testing.T) {
	defer destroyAfterTest()
	r := require.New(t)
	// TestListenToSamePort 测试监听同一个端口
	_, err := Listen(TcpNet, "127.0.0.1:"+TestPort1, agent.GetEchoAgent, core.WithMaxConnNum(TestThread))
	r.Nil(err)
	_, err = Listen(TcpNet, "127.0.0.1:"+TestPort1, agent.GetEchoAgent, core.WithMaxConnNum(TestThread))
	r.NotNil(err)
}

// 测试使用不支持的协议
func TestListenToInvalidProtocol(t *testing.T) {
	defer destroyAfterTest()
	_, err := Listen(UdpSeries-1, "127.0.0.1:"+TestPort1, agent.GetEchoAgent, core.WithMaxConnNum(TestThread))
	assert.NotNil(t, err)
	DestroyAll()
	time.Sleep(DestroyAllowWaitTime)
}

// 测试使用不合法的 Agent
func TestInvalidAgent(t *testing.T) {
	defer destroyAfterTest()
	r := require.New(t)
	netList := GetInfo()
	for typ, i := range netList {
		t.Logf("test network type: %d", typ)
		if i.SupportServer() {
			_, err := Listen(typ, "127.0.0.1:"+TestPort1, nil, core.WithMaxConnNum(TestThread))
			r.Equal(err, core.ErrInvalidGetAgentFunc)
		}
		if i.SupportClient() {
			_, err := Connect(typ, "127.0.0.1:"+TestPort1, nil, core.WithReconnect(true))
			r.Equal(err, core.ErrInvalidGetAgentFunc)
		}
	}
}

// 测试使用不合法的 地址
func TestInvalidAddress(t *testing.T) {
	defer destroyAfterTest()
	r := require.New(t)
	netList := GetInfo()
	for typ, i := range netList {
		t.Logf("test network type: %d", typ)
		if i.SupportServer() {
			_, err := Listen(typ, ":-1", agent.GetSingleAgent, core.WithMaxConnNum(TestThread))
			r.Equal(err, core.ErrInvalidAddress)
		}
		if i.SupportClient() {
			_, err := Connect(typ, ":-1", agent.GetSingleAgent, core.WithReconnect(true))
			r.Equal(err, core.ErrInvalidAddress)
		}
	}
}

// 测试正常连接和收发消息
func TestNetworkNormal(t *testing.T) {
	defer destroyAfterTest()
	r := require.New(t)
	netList := GetInfo()
	for typ, inf := range netList {
		t.Logf("test network type: %d", typ)
		count = 0

		_, err := Listen(typ, "127.0.0.1:"+TestPort1, agent.GetEchoAgent, core.WithMaxConnNum(TestThread))
		time.Sleep(ListenAllowWaitTime)
		r.Nil(err)
		for i := 0; i < TestThread+TestWaitThread; i++ {
			if inf.SupportClient() {
				_, err = Connect(typ, "127.0.0.1:"+TestPort1, newTestSendAgent, core.WithReconnect(true))
			} else {
				_, err = Connect(TcpNet, "127.0.0.1:"+TestPort1, newTestSendAgent, core.WithReconnect(true))
			}
			r.Nil(err)
		}

		time.Sleep(WaitMsgSendTime)
		r.Equal(SendMsgNum*TestThread, count)
		r.Equal(TestThread*2+TestWaitThread, GetConnNum())
		destroyAfterTest()

		r.Equal(0, GetConnNum())
	}
}

// 测试客户端主动断开连接

// 测试服务端取消监听，然后再开启监听
