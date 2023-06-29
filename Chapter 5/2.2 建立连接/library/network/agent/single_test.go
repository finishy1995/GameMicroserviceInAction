package agent

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSingleAgent(t *testing.T) {
	conn := new(mockConn)
	r := require.New(t)
	var a string
	conn.handle = func(b []byte) (n int, err error) {
		a = string(b)
		return len(a), errors.New("test error")
	}
	agent := GetSingleAgent()
	r.NotNil(agent)

	agent.OnConnect(conn)
	agent.OnMessage([]byte(""), conn)
	agent.OnClose(conn)
}
