package agent

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendAgent_OnConnect(t *testing.T) {
	conn := new(mockConn)
	r := require.New(t)
	var a string
	conn.handle = func(b []byte) (n int, err error) {
		a = string(b)
		return len(a), errors.New("test error")
	}
	agent := GetSendAgent()
	r.NotNil(agent)

	agent.OnConnect(conn)
	r.Equal("hello world", a)
}
