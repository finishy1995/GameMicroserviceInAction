// Package bytebuffer 取代原生 bytes.Buffer，性能更高且更少的内存碎片
//
// Example Usage
//
//	bb := bytebuffer.Get()
//	defer bytebuffer.Put(bb)
//	bb.Write([]byte("111"))
//	bb.Write([]byte("222"))
//	bb.Bytes() == "111222"
//
// Benchmark History
//
//	goos: windows
//	goarch: amd64
//	pkg: nbserver/common/buffer/bytebuffer
//	cpu: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz
//	BenchmarkGetWrite（Get，碎片更少）
//	BenchmarkGetWrite-16            27677829                44.18 ns/op
//	BenchmarkNewWrite（New 一个新对象不复用）
//	BenchmarkNewWrite-16            31481107                38.80 ns/op
//	BenchmarkOriginWrite（系统自带）
//	BenchmarkOriginWrite-16         19354963                62.00 ns/op
package bytebuffer

import "github.com/valyala/bytebufferpool"

type ByteBuffer = bytebufferpool.ByteBuffer

var (
	// Get 从池子中获取一个空的 byte buffer
	Get = bytebufferpool.Get

	// Put 返回一个 byte buffer 回到池子中
	Put = func(b *ByteBuffer) {
		if b != nil {
			bytebufferpool.Put(b)
		}
	}
)
