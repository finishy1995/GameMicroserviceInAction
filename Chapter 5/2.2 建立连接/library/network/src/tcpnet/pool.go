package tcpnet

import "sync"

// TODO CPU 缓存优化。实现一个pool提前分配好一定内存的 slice，这样循环的时候，l1缓存能够一次存储更多的 connection
var pool = sync.Pool{
	New: func() interface{} {
		return new(Conn)
	},
}
