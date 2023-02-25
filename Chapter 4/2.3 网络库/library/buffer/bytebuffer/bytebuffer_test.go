package bytebuffer

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestByteBuffer(t *testing.T) {
	r := require.New(t)
	bb := Get()
	defer Put(bb)
	_, err := bb.WriteString("hello")
	r.Nil(err)
	err = bb.WriteByte(' ')
	r.Nil(err)
	_, err = bb.Write([]byte("world"))
	r.Nil(err)
	r.Equal("hello world", bb.String())
}

func BenchmarkGetWrite(b *testing.B) {
	s := []byte("foobarbaz")
	b.RunParallel(func(pb *testing.PB) {
		buf := Get()
		for pb.Next() {
			for i := 0; i < 100; i++ {
				buf.Write(s)
			}
			buf.Reset()
		}
	})
}

func BenchmarkNewWrite(b *testing.B) {
	s := []byte("foobarbaz")
	b.RunParallel(func(pb *testing.PB) {
		var buf ByteBuffer
		for pb.Next() {
			for i := 0; i < 100; i++ {
				buf.Write(s)
			}
			buf.Reset()
		}
	})
}

func BenchmarkOriginWrite(b *testing.B) {
	s := []byte("foobarbaz")
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			for i := 0; i < 100; i++ {
				buf.Write(s)
			}
			buf.Reset()
		}
	})
}
