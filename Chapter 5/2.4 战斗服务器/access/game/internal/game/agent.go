package game

import (
	"ProjectX/access/game/pb/game"
	"ProjectX/library/log"
	agentBase "ProjectX/library/network/agent"
	"ProjectX/library/network/core"
	"github.com/golang/protobuf/proto"
)

type agent struct {
	agentBase.SingleAgent
	game   *Game
	conn   core.Conn
	userId string
}

func newAgent(game *Game) core.Agent {
	if game == nil {
		panic("game is nil")
	}
	return &agent{
		game: game,
	}
}

func (a *agent) OnConnect(conn core.Conn) {
	a.conn = conn
}

func (a *agent) OnMessage(b []byte, _ core.Conn) {
	req := &game.Action{}

	if err := proto.Unmarshal(b, req); err != nil {
		log.Error("unmarshal error: %v", err)
		return
	}

	a.game.OnRpc(a, req)
}

func (a *agent) send(b []byte) {
	_, err := a.conn.Write(b)
	if err != nil {
		log.Error("write error: %v", err)
		return
	}
}
