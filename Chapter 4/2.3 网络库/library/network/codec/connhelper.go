package codec

import (
	"ProjectX/library/buffer/bytebuffer"
	"ProjectX/library/buffer/ringbuffer"
)

// ConnHelper Codec 连接帮助
type ConnHelper struct {
	inboundBuffer *ringbuffer.RingBuffer
	buffer        *bytebuffer.ByteBuffer
	byteBuffer    *bytebuffer.ByteBuffer
}

// InitBuffer 初始化 buffer
func (conn *ConnHelper) InitBuffer() {
	conn.inboundBuffer = ringbuffer.Get()
	conn.buffer = bytebuffer.Get()
}

// PushPacket 把网络包压入 buffer
func (conn *ConnHelper) PushPacket(b []byte) {
	conn.inboundBuffer.Write(b)
}

// Read 读取所有数据，不移动读指针
func (conn *ConnHelper) Read() (buf []byte) {
	if conn.inboundBuffer.IsEmpty() {
		return conn.buffer.Bytes()
	}
	conn.byteBuffer = conn.inboundBuffer.WithByteBuffer(conn.buffer.Bytes())
	return conn.byteBuffer.Bytes()
}

// ResetBuffer 重置读取容器
func (conn *ConnHelper) ResetBuffer() {
	conn.buffer.Reset()
	conn.inboundBuffer.Reset()
	bytebuffer.Put(conn.byteBuffer)
	conn.byteBuffer = nil
}

// ReadN 读取给定长度的数据，如果数据不够，则返回所有数据，不移动读指针
func (conn *ConnHelper) ReadN(n int) (size int, buf []byte) {
	inBufferLen := conn.inboundBuffer.Length()
	tempBufferLen := conn.buffer.Len()
	if totalLen := inBufferLen + tempBufferLen; totalLen < n || n <= 0 {
		n = totalLen
	}
	size = n
	if conn.inboundBuffer.IsEmpty() {
		buf = conn.buffer.B[:n]
		return
	}
	head, tail := conn.inboundBuffer.LazyRead(n)
	conn.byteBuffer = bytebuffer.Get()
	_, _ = conn.byteBuffer.Write(head)
	_, _ = conn.byteBuffer.Write(tail)
	if inBufferLen >= n {
		buf = conn.byteBuffer.Bytes()
		return
	}

	restSize := n - inBufferLen
	_, _ = conn.byteBuffer.Write(conn.buffer.B[:restSize])
	buf = conn.byteBuffer.Bytes()
	return
}

// ShiftN 移动读指针到给定长度
func (conn *ConnHelper) ShiftN(n int) (size int) {
	inBufferLen := conn.inboundBuffer.Length()
	tempBufferLen := conn.buffer.Len()
	if inBufferLen+tempBufferLen < n || n <= 0 {
		conn.ResetBuffer()
		size = inBufferLen + tempBufferLen
		return
	}
	size = n
	if conn.inboundBuffer.IsEmpty() {
		conn.buffer.B = conn.buffer.B[n:]
		return
	}

	bytebuffer.Put(conn.byteBuffer)
	conn.byteBuffer = nil

	if inBufferLen >= n {
		conn.inboundBuffer.Shift(n)
		return
	}
	conn.inboundBuffer.Reset()

	restSize := n - inBufferLen
	conn.buffer.B = conn.buffer.B[restSize:]
	return
}

// BufferLength 读取容器数据长度
func (conn *ConnHelper) BufferLength() (size int) {
	return conn.inboundBuffer.Length() + conn.buffer.Len()
}
