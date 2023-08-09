package network

import "ProjectX/library/network/core"

// NetType 网络类型
type NetType int16

type serverMaker func() core.Server

type clientMaker func() core.Client

// info 信息
type info struct {
	// client 客户端
	client clientMaker
	// server 服务端
	server serverMaker
	// codecSupport Codec 支持
	codecSupport bool
}

const (
	// TcpSeries 从 0 开始，TCP系列不同的实现
	TcpSeries NetType = iota
	// UdpSeries UDP 系列的不同实现
	UdpSeries = TcpSeries + SeriesInterval

	// SeriesInterval 系列间隔
	SeriesInterval NetType = 1000

	TcpNet       = TcpSeries + 0
	TcpGNet      = TcpSeries + 1
	WebsocketNet = TcpSeries + 2
)

// SupportClient 支持客户端
func (i *info) SupportClient() bool {
	return i.client != nil
}

// SupportServer 支持服务端
func (i *info) SupportServer() bool {
	return i.server != nil
}
