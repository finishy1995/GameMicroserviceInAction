package agent

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

type mockConn struct {
	handle func(_ []byte) (n int, err error)
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

func TestEchoAgent_OnMessage(t *testing.T) {
	conn := new(mockConn)
	r := require.New(t)
	var a string
	conn.handle = func(b []byte) (n int, err error) {
		a = string(b)
		return len(a), errors.New("test error")
	}
	agent := GetEchoAgent()
	r.NotNil(agent)
	b := []byte("hello world")

	agent.OnMessage(b, conn)
	r.Equal("hello world", a)
}
