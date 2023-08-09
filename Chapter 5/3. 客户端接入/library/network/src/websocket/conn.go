package websocket

import (
	"ProjectX/library/log"
	"ProjectX/library/network/core"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

type Conn struct {
	*websocket.Conn
	agent core.Agent
}

func (c *Conn) Run() {
}

func (c *Conn) Close() {
	err := c.Conn.Close()
	if err != nil {
		log.Debug("close websocket conn error: %v", err)
	}
	c.agent.OnClose(c)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	err = c.Conn.WriteMessage(websocket.TextMessage, b)
	return len(b), err
}

func newConn(webConn *websocket.Conn, agent core.Agent) *Conn {
	return &Conn{
		Conn:  webConn,
		agent: agent,
	}
}
