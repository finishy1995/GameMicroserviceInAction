package ringbuffer

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	THREAD = 200
)

func TestRingBuffer(t *testing.T) {
	r := require.New(t)
	var buf RingBuffer
	buf.WriteString("hello")
	// 不改变读指针位置
	head, _ := buf.LazyRead(5)
	s := string(head[:])
	r.Equal("hello", s)
	b := make([]byte, 10)
	// 改变读指针位置
	n, err := buf.Read(b)
	r.Nil(err)
	r.Equal(5, n)
	s = string(b[:n])
	r.Equal("hello", s)
}

// 测试并行情况下的线程安全(一个协程写，一个协程读，不支持多协程读写)
func TestRingBufferConcurrent(t *testing.T) {
	var buf RingBuffer
	// 测试并发读和数据一致性
	check := make(map[byte]bool)
	for i := 0; i < THREAD; i++ {
		a := byte(i)
		check[a] = false
	}

	go func() {
		for i := 0; i < THREAD; i++ {
			a := byte(i)
			buf.WriteByte(a)
		}
	}()
	go func() {
		count := 0
		for {
			p := make([]byte, 1)
			n, _ := buf.Read(p)
			if n == 1 {
				check[p[0]] = true
				count++
			}
			if count == THREAD {
				break
			}
		}
	}()

	// 保证并发执行完成
	time.Sleep(time.Millisecond * 20)
	count := 0
	for _, value := range check {
		if value {
			count++
		}
	}
	assert.Equal(t, THREAD, count)
}

func TestGetPut(t *testing.T) {
	buf := Get()
	buf.WriteByte('a')
	Put(buf)
	Put(nil)
}
