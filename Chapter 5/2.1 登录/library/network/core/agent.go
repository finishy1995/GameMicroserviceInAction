package core

// GetAgent 获取代理
type GetAgent func() Agent

// Agent 网络代理
type Agent interface {
	// OnConnect 连接创建
	OnConnect(conn Conn)
	// OnMessage 收到消息
	OnMessage(b []byte, conn Conn)
	// OnClose 连接关闭
	OnClose(conn Conn)
}
