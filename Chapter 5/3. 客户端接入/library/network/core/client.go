package core

import (
	"time"
)

const (
	DefaultConnectWaitStart = time.Millisecond * 20
	DefaultConnectWaitMut   = 2
	DefaultConnectMaxWait   = time.Second * 2
)

// Client 网络客户端
type Client interface {
	Object
	// Start 开启客户端连接
	Start(address string, newAgent GetAgent, opts ...ClientOption) error
	// IsConnected 是否处于连接状态
	IsConnected() bool
}

// ClientOption 客户端配置项
type ClientOption func(*ClientOptions)

// ClientOptions 配置结构体
type ClientOptions struct {
	// 是否自动重连 默认为否
	Reconnect bool
	// 特定客户端参数
	Context map[string]interface{}
}

// WithReconnect 重连配置
func WithReconnect(reconnect bool) ClientOption {
	return func(o *ClientOptions) {
		o.Reconnect = reconnect
	}
}

// WithClientContext 特定参数配置
func WithClientContext(context map[string]interface{}) ClientOption {
	return func(o *ClientOptions) {
		o.Context = context
	}
}

var (
	// DefaultClientOptions 默认 Client 选项
	DefaultClientOptions = ClientOptions{
		Reconnect: false,
		Context:   nil,
	}
)
