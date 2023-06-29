package codec

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

type mockConn struct {
	ConnHelper
	handle func(_ []byte) (n int, err error)
}

func newMockConn() *mockConn {
	c := new(mockConn)
	c.InitBuffer()
	return c
}

func (conn *mockConn) Run()                 {}
func (conn *mockConn) Close()               {}
func (conn *mockConn) LocalAddr() net.Addr  { return nil }
func (conn *mockConn) RemoteAddr() net.Addr { return nil }
func (conn *mockConn) Write(b []byte) (n int, err error) {
	if conn.handle != nil {
		return conn.handle(b)
	}
	return 0, nil
}
func (conn *mockConn) MockGetNetworkMsg(b []byte) {
	conn.PushPacket(b)
}

func TestBuiltInCodec(t *testing.T) {
	codec := new(BuiltInCodec)
	conn := newMockConn()
	r := require.New(t)

	test1 := []byte("hello world")
	bb, err := codec.Encode(conn, test1)
	r.Nil(err)
	r.Equal(test1, bb)
	conn.MockGetNetworkMsg(bb)
	bb, err = codec.Decode(conn)
	r.Nil(err)
	r.Equal(test1, bb)

	test2 := []byte("")
	bb, err = codec.Encode(conn, test2)
	r.Nil(err)
	r.Equal(test2, bb)
	conn.MockGetNetworkMsg(bb)
	bb, err = codec.Decode(conn)
	r.Nil(err)
	r.Equal(test2, bb)
}
