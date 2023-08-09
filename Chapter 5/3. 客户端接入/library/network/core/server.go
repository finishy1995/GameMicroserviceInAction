package core

import (
	"time"
)

const (
	DefaultMaxConnNum int = 3000
	DefaultKeepAlive      = time.Second * 1
	UpdateInterval        = time.Millisecond * 10
)

// Server 网络服务器
type Server interface {
	Object
	// Start 开启服务器
	Start(address string, newAgent GetAgent, opts ...ServerOption) error
	// GetConnNum 获取服务器连接数
	GetConnNum() (num int)
}

// ServerOption 服务器配置项
type ServerOption func(*ServerOptions)

// ServerOptions 服务器配置结构体
type ServerOptions struct {
	// 最大连接数 默认为-1（不限）
	MaxConnNum int
	// 特定服务器参数
	Context map[string]interface{}
}

// WithMaxConnNum 最大连接数配置
func WithMaxConnNum(maxConnNum int) ServerOption {
	return func(o *ServerOptions) {
		o.MaxConnNum = maxConnNum
	}
}

// WithServerContext 特定参数配置
func WithServerContext(context map[string]interface{}) ServerOption {
	return func(o *ServerOptions) {
		o.Context = context
	}
}

var (
	// DefaultServerOptions 默认 Server 选项
	DefaultServerOptions = ServerOptions{
		MaxConnNum: -1,
		Context:    nil,
	}
)
