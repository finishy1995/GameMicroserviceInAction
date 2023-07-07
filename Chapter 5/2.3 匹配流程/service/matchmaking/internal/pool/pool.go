package pool

import (
	"ProjectX/library/routine"
	"fmt"
	"sync"
	"time"
)

// Pool 匹配池
type Pool struct {
	ticketIdInfoMapping map[string]*info  // ticketId -> info
	userIdTicketMapping map[string]string // userId -> ticketId
	preparingMapping    map[string]*info  // ticketId -> info
	preparingMutex      sync.Mutex        // preparingMapping 互斥锁
	closeSig            chan bool         // 关闭信号
}

func NewPool() *Pool {
	p := &Pool{
		ticketIdInfoMapping: make(map[string]*info),
		userIdTicketMapping: make(map[string]string),
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
	// 将 preparingMapping 中候补的人员加入匹配池 ticketIdInfoMapping
	p.preparingMutex.Lock()
	for ticketId, info := range p.preparingMapping {
		if info.status != StatusPreparing {
			continue
		}

		info.status = StatusMatching
		p.ticketIdInfoMapping[ticketId] = info
	}
	p.preparingMapping = make(map[string]*info, 0)
	p.preparingMutex.Unlock()

	// TODO: 匹配逻辑，根据玩家信息按需要匹配
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
	p.preparingMutex.Lock()
	defer p.preparingMutex.Unlock()

	// 如果用户已经在匹配池中，直接返回 ticketId
	if ticketId, ok := p.userIdTicketMapping[userId]; ok {
		return ticketId
	}

	// 创建匹配信息
	ticketId := fmt.Sprintf("%s_%d", userId, time.Now().Unix())
	info := &info{
		ticketId:    ticketId,
		createdTime: time.Now().Unix(),
		status:      StatusPreparing,
	}
	p.preparingMapping[ticketId] = info
	p.userIdTicketMapping[userId] = ticketId
	return ticketId
}

// RemoveUser 从匹配池中移除用户
func (p *Pool) RemoveUser(userId string) {
	if ticketId, ok := p.userIdTicketMapping[userId]; ok {
		p.preparingMutex.Lock()
		defer p.preparingMutex.Unlock()

	} else {
		return
	}
}
