package agent

import "ProjectX/library/network/core"

// SingleAgent 最简单的网络代理，单例且不响应任何事件
type SingleAgent struct {
}

var (
	singleInstance = new(SingleAgent)
)

func GetSingleAgent() core.Agent {
	return singleInstance
}

func (agent *SingleAgent) OnConnect(_ core.Conn) {
}

func (agent *SingleAgent) OnMessage(_ []byte, _ core.Conn) {
}

func (agent *SingleAgent) OnClose(_ core.Conn) {
}
