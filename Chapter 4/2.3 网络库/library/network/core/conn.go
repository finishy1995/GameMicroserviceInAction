package core

import (
	"net"
)

const (
	TCPMaxPackageSize       = 10240
	TimeoutTime       int64 = 7000
	HeartbeatTime     int64 = 2000
)

// Conn 网络连接
type Conn interface {
	Object
	// Write 写入并发送数据
	Write(b []byte) (n int, err error)
	// LocalAddr 本地地址
	LocalAddr() net.Addr
	// RemoteAddr 远程地址
	RemoteAddr() net.Addr
}

// CodecConn 支持 Codec 的网络连接
type CodecConn interface {
	Conn

	// Read 读取所有数据，不移动读指针
	Read() (buf []byte)

	// ResetBuffer 重置读取容器
	ResetBuffer()

	// ReadN 读取给定长度的数据，如果数据不够，则返回所有数据，不移动读指针
	ReadN(n int) (size int, buf []byte)

	// ShiftN 移动读指针到给定长度
	ShiftN(n int) (size int)

	// BufferLength 读取容器数据长度
	BufferLength() (size int)
}
