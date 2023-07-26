package game

import (
	"ProjectX/access/game/pb/game"
	"ProjectX/library/log"
	"ProjectX/library/network"
	"ProjectX/library/network/core"
	"ProjectX/library/routine"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
)

type Game struct {
	id       uint32
	port     uint32
	frame    int32
	status   Status
	secret   string
	player   map[string]*agent
	deadline time.Time
}

func NewGame(port uint32) *Game {
	g := &Game{
		status: Idle,
		id:     port,
		port:   port,
	}
	err := routine.Run(true, g.run)
	if err != nil {
		panic(err)
	}

	_, err = network.Listen(network.TcpGNet, fmt.Sprintf("0.0.0.0:%d", port), func() core.Agent {
		return newAgent(g)
	}, core.WithMaxConnNum(100))
	if err != nil {
		panic(err)
	}

	return g
}

func (g *Game) run() {
	ticker := time.NewTicker(time.Second / 25) // 25 帧
	for {
		select {
		case <-ticker.C:
			g.update()
		}
	}
}

func (g *Game) update() {
	g.frame++
	switch g.status {
	case Lobby:
		if time.Now().After(g.deadline) {
			g.transfer2Status(BeforeStart)
		}
		return
	case BeforeStart:
		if time.Now().After(g.deadline) {
			g.transfer2Status(Playing)
		}
		return
	case Playing:
		if time.Now().After(g.deadline) {
			g.transfer2Status(BeforeTerminate)
		}
		return
	case BeforeTerminate:
		if time.Now().After(g.deadline) {
			g.status = Idle
		}
		return
	}
}

func (g *Game) Status() Status {
	return g.status
}

func (g *Game) GetPort() uint32 {
	return g.port
}

func (g *Game) Init(secret string, player []string) {
	g.secret = secret
	g.player = make(map[string]*agent, len(player))
	for _, p := range player {
		g.player[p] = nil
	}
	g.transfer2Status(Lobby)
}

func (g *Game) AddPlayer(secret string, player string, a *agent) {
	if g.status != Lobby {
		return // TODO：战局内中途加入/掉线重连
	}
	if g.secret != secret {
		return
	}
	if _, ok := g.player[player]; !ok {
		return
	}
	flag := true
	b, err := proto.Marshal(&game.Action{
		ActionType:    PlayerJoinAction,
		ActionParam:   fmt.Sprintf(`{"player":"%s"}`, player),
		ActionFrameId: g.frame,
	})
	if err != nil {
		log.Error("marshal error: %v", err)
		return
	}

	g.player[player] = a
	for _, v := range g.player {
		if v == nil {
			flag = false
		}
		v.send(b)
	}
	if flag {
		g.transfer2Status(BeforeStart)
	}
}

func (g *Game) RemovePlayer(player string) {
	if _, ok := g.player[player]; !ok {
		return
	}
	g.player[player] = nil
	g.broadcast(&game.Action{
		ActionType:    PlayerLeaveAction,
		ActionParam:   fmt.Sprintf(`{"player":"%s"}`, player),
		ActionFrameId: g.frame,
	})
}

func (g *Game) transfer2Status(status Status) {
	if g.status == status {
		return
	}
	g.status = status
	switch status {
	case Lobby:
		g.frame = 0
		g.deadline = time.Now().Add(60 * time.Second) // 大厅倒计时
		return
	case BeforeStart:
		g.deadline = time.Now().Add(5 * time.Second) // 开始倒计时
		g.broadcast(&game.Action{
			ActionType:    StartAction,
			ActionParam:   fmt.Sprintf(`{"start":%d}`, g.deadline.UnixMilli()),
			ActionFrameId: g.frame,
		})
		return
	case Playing:
		g.deadline = time.Now().Add(30 * time.Minute) // 一局最多 30 分钟
		return
	case BeforeTerminate:
		g.deadline = time.Now().Add(20 * time.Second) // 结束倒计时
		g.broadcast(&game.Action{
			ActionType:    TerminateAction,
			ActionParam:   fmt.Sprintf(`{"terminate":%d}`, g.deadline.UnixMilli()),
			ActionFrameId: g.frame,
		})
		return
	}
}

func (g *Game) OnRpc(a *agent, action *game.Action) {
	if a == nil || action == nil {
		return
	}
	if a.userId == "" && action.ActionType != InitAction {
		return
	}
	var context map[string]interface{}
	err := json.Unmarshal([]byte(action.ActionParam), &context)
	if err != nil {
		return
	}

	switch action.ActionType {
	case InitAction:
		g.AddPlayer(context["secret"].(string), context["id"].(string), a)
		return
	default:
		g.broadcast(action)
		return
	}
}

func (g *Game) broadcast(action *game.Action) {
	if action.ActionFrameId == 0 {
		action.ActionFrameId = g.frame
	}
	b, err := proto.Marshal(action)
	if err != nil {
		log.Error("marshal error: %v", err)
		return
	}
	for _, v := range g.player {
		if v == nil {
			continue
		}
		v.send(b)
	}
}
