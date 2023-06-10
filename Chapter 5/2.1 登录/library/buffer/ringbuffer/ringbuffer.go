// Package ringbuffer ring buffer 无锁循环缓冲区（线程安全），首尾相接的循环数组，详见算法
//
// Example Usage
//
//	rb := ringbuffer.Get()
//	defer ringbuffer.Put(rb)
//	rb.WriteString("hello")
//	b := make([]byte, 10)
//	// 改变读指针位置，不改变使用 LazyRead
//	n, err := rb.Read(b) // n = 5; b = []byte("hello")
package ringbuffer

import (
	ringbuffer3rd "github.com/panjf2000/gnet/pool/ringbuffer"
)

type RingBuffer = ringbuffer3rd.RingBuffer

var (
	// Get 从池子中获取一个空的 ring buffer
	Get = ringbuffer3rd.Get

	// Put 返回一个 ring buffer 回到池子中
	Put = func(b *RingBuffer) {
		if b != nil {
			ringbuffer3rd.Put(b)
		}
	}
)
