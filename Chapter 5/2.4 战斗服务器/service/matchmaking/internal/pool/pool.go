package pool

import (
	"ProjectX/library/routine"
	"fmt"
	"sync"
	"time"
)

type action struct {
	userId   string
	ticketId string
	action   Action
}

// Pool 匹配池
type Pool struct {
	sync.RWMutex
	ticketIdInfoMapping map[string]*Info  // ticketId -> info 正在匹配中的 ticket
	userIdTicketMapping map[string]string // userId -> ticketId 用户与 ticket 的映射
	backupMapping       map[string]*Info  // ticketId -> info 备份的 ticket
	preparingAction     chan *action      // 准备中的动作
	closeSig            chan bool         // 关闭信号
}

func NewPool() *Pool {
	p := &Pool{
		ticketIdInfoMapping: make(map[string]*Info, 0),
		userIdTicketMapping: make(map[string]string, 0),
		backupMapping:       make(map[string]*Info, 0),
		preparingAction:     make(chan *action, 1000),
		closeSig:            make(chan bool, 2),
	}
	err := routine.Run(false, p.run)
	if err != nil {
		panic(err)
	}
	return p
}

func (p *Pool) run() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			p.match()
			break
		case <-p.closeSig:
			return
		}
	}
}

func (p *Pool) match() {
	p.Lock()
	// 维护 backupMapping，将其中超时的人员删除
	now := time.Now().Unix()
	for ticketId, info := range p.backupMapping {
		if info.endTime == 0 {
			continue
		}
		if now-info.endTime > 900 {
			delete(p.backupMapping, ticketId)
			delete(p.userIdTicketMapping, info.userId)
		}
	}

	// 将 preparingAction 中的动作一一执行
	for a := range p.preparingAction {
		switch a.action {
		case ActionAdd:
			p.addUser(a.userId, a.ticketId)
			break
		case ActionCancel:
			p.cancelUser(a.userId)
			break
		}
	}
	p.Unlock()

	// TODO: 匹配逻辑，根据玩家信息按需要匹配；添加超时逻辑
	// 临时逻辑是，只要有两个人匹配，就完成匹配
	for {
		if len(p.ticketIdInfoMapping) < 2 {
			break
		}
		var ticketIds []string
		for ticketId := range p.ticketIdInfoMapping {
			ticketIds = append(ticketIds, ticketId)
			if len(ticketIds) == 2 {
				break
			}
		}
		if len(ticketIds) == 2 {
			p.matchSuccess(ticketIds)
		}
	}
}

func (p *Pool) matchSuccess(ids []string) {
	for _, ticketId := range ids {
		info := p.ticketIdInfoMapping[ticketId]
		info.status = StatusMatched
		info.endTime = time.Now().Unix()
		delete(p.ticketIdInfoMapping, ticketId)
	}
}

// Close 关闭匹配池
func (p *Pool) Close() {
	p.closeSig <- true
}

// AddUser 添加用户到匹配池
func (p *Pool) AddUser(userId string) string {
	// 如果用户已经在匹配池中，直接返回 ticketId
	p.RLock()
	if ticketId, ok := p.userIdTicketMapping[userId]; ok {
		if p.backupMapping[ticketId].status == StatusPreparing || p.backupMapping[ticketId].status == StatusMatching {
			return ticketId
		}
	}
	p.RUnlock()

	// 创建匹配信息
	ticketId := fmt.Sprintf("%s_%d", userId, time.Now().Unix())

	p.preparingAction <- &action{
		userId:   userId,
		ticketId: ticketId,
		action:   ActionAdd,
	}

	return ticketId
}

// CancelUser 从匹配池中移除用户
func (p *Pool) CancelUser(userId string) {
	p.preparingAction <- &action{
		userId: userId,
		action: ActionCancel,
	}
}

func (p *Pool) GetInfo(userId string) *Info {
	p.RLock()
	defer p.RUnlock()
	if ticketId, ok := p.userIdTicketMapping[userId]; ok {
		if info, ok := p.backupMapping[ticketId]; ok {
			return info
		}
	}
	return nil
}

func (p *Pool) addUser(userId string, ticketId string) {
	info := &Info{
		ticketId:    ticketId,
		userId:      userId,
		createdTime: time.Now().Unix(),
		status:      StatusPreparing,
	}
	p.ticketIdInfoMapping[ticketId] = info
	p.userIdTicketMapping[userId] = ticketId
	p.backupMapping[ticketId] = info
}

func (p *Pool) cancelUser(userId string) {
	if ticketId, ok := p.userIdTicketMapping[userId]; ok {
		if info, ok := p.backupMapping[ticketId]; ok {
			info.status = StatusCanceled
			info.endTime = time.Now().Unix()
		}
		delete(p.ticketIdInfoMapping, ticketId)
	}
}
