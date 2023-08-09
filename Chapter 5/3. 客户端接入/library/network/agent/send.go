package agent

import (
	"ProjectX/library/log"
	"ProjectX/library/network/core"
)

// SendAgent 发消息代理，一连上就开始发消息
type SendAgent struct {
	SingleAgent
}

func GetSendAgent() core.Agent {
	return new(SendAgent)
}

func (agent *SendAgent) OnConnect(conn core.Conn) {
	_, err := conn.Write([]byte("hello world"))
	if err != nil {
		log.Error("connection write msg error: %s", err.Error())
	}
}
