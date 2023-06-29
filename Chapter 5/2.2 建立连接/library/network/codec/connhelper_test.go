package codec

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	helper *ConnHelper = nil
)

func initOnce() {
	if helper != nil {
		return
	}
	helper = new(ConnHelper)
	helper.InitBuffer()
}

func TestConnHelper_InitBuffer(t *testing.T) {
	initOnce()
}

func TestConnHelper_Read(t *testing.T) {
	initOnce()
	a := helper.Read()
	r := require.New(t)
	r.Empty(a)
	helper.PushPacket([]byte("hello world"))
	a = helper.Read()
	r.Equal(a, []byte("hello world"))
	r.Equal(11, helper.BufferLength())
	helper.ResetBuffer()
	a = helper.Read()
	r.Empty(a)
}

func TestConnHelper_ReadN(t *testing.T) {
	initOnce()
	r := require.New(t)
	helper.PushPacket([]byte("hello world"))
	size, a := helper.ReadN(0)
	r.Equal(11, size)
	r.Equal([]byte("hello world"), a)
	size, a = helper.ReadN(5)
	r.Equal([]byte("hello"), a)
	size, a = helper.ReadN(16)
	r.Equal(11, size)
	r.Equal([]byte("hello world"), a)

	size = helper.ShiftN(3)
	r.Equal(3, size)
	size, a = helper.ReadN(0)
	r.Equal(8, size)
	r.Equal([]byte("lo world"), a)
	size = helper.ShiftN(16)
	r.Equal(8, size)
	size, a = helper.ReadN(0)
	r.Equal(0, size)
	r.Empty(a)
	helper.PushPacket([]byte("test"))
	r.Equal(4, helper.BufferLength())
	size = helper.ShiftN(0)
	r.Equal(4, size)
	r.Equal(0, helper.BufferLength())
}
