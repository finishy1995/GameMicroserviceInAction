package agent

import (
	"ProjectX/library/log"
	"ProjectX/library/network/core"
)

// EchoAgent 最简单的 Echo 代理，发什么回什么，使用单例实现（也可以和 Conn 一一对应）
type EchoAgent struct {
	SingleAgent
}

var (
	echoInstance = new(EchoAgent)
)

func GetEchoAgent() core.Agent {
	return echoInstance
}

func (agent *EchoAgent) OnMessage(b []byte, conn core.Conn) {
	_, err := conn.Write(b)
	if err != nil {
		log.Error("connection write msg error: %s", err.Error())
	}
}
