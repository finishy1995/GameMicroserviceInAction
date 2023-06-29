package agent

import (
	"ProjectX/access/gate/pb/gate"
	"ProjectX/base"
	"ProjectX/library/log"
	"ProjectX/library/network/core"
	"ProjectX/library/routine"
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/timex"
	"strings"
	"sync"
	"time"
)

type Agent struct {
	mutex     sync.Mutex
	userId    string
	rpcTask   chan *task
	conn      core.Conn
	closeSig  chan bool
	closeFlag bool
}

const (
	MaxChanSize   = 100
	RpcTimeout    = 5 * time.Second
	SlowThreshold = 40 * time.Millisecond
)

var (
	userAgentMapping      = make(map[string]*Agent)
	userAgentMappingMutex sync.Mutex
)

func NewAgent() core.Agent {
	agt := &Agent{
		rpcTask:   make(chan *task, MaxChanSize),
		closeFlag: false,
		closeSig:  make(chan bool, 2), // 这里设置为2，防止重复关闭的阻塞，充分 debug 后可设置为 1
	}
	err := routine.Run(true, agt.run)
	if err != nil {
		panic(err)
	}
	return agt
}

func GetUserNum() int {
	return len(userAgentMapping)
}

// OnConnect is called when a connection is first established.
func (agt *Agent) OnConnect(conn core.Conn) {
	agt.conn = conn

	// TODO: 这里可以做一些初始化操作，例如获取玩家 session 信息、获取玩家匹配状态对局信息等
}

// OnMessage is called when a new message is received from the connection.
func (agt *Agent) OnMessage(b []byte, conn core.Conn) {
	agt.Send(agt.handleMessage(b), conn)
}

// OnClose is called when the connection closed.
func (agt *Agent) OnClose(_ core.Conn) {
	// TODO: 这里可以做一些清理操作，例如玩家匹配状态对局信息清理等

	if agt.userId != "" {
		userAgentMappingMutex.Lock()
		defer userAgentMappingMutex.Unlock()
		delete(userAgentMapping, agt.userId)
	}

	agt.closeFlag = true
	close(agt.closeSig)
}

func (agt *Agent) handleMessage(b []byte) *gate.ClientResponse {
	if agt.userId == "" {
		agt.mutex.Lock()
		defer agt.mutex.Unlock()
		if agt.userId == "" {
			return agt.handleToken(b)
		}
	}

	return agt.handleMethod(b)
}

func (agt *Agent) Send(resp *gate.ClientResponse, conn core.Conn) {
	if resp == nil {
		return
	}

	result, err := proto.Marshal(resp)
	if err != nil {
		log.Error("userId : %s, error : %s", agt.userId, err.Error())
	} else {
		_, err = conn.Write(result)
		if err != nil {
			log.Error("userId : %s, msg write error : %s", agt.userId, err.Error())
		}
	}
}

func (agt *Agent) handleToken(b []byte) *gate.ClientResponse {
	resp := &gate.ClientResponse{
		Code: base.ErrorCodeInvalidGateToken,
	}

	// TODO: 使用加密算法解析 token，获取 userId
	userId := string(b)
	if userId == "" {
		return resp
	}

	agt.userId = userId
	resp.Code = base.ErrorCodeOK
	userAgentMappingMutex.Lock()
	defer userAgentMappingMutex.Unlock()
	userAgentMapping[userId] = agt
	return resp
}

type task struct {
	method  string
	content []byte
	id      int32
}

func (agt *Agent) handleMethod(b []byte) *gate.ClientResponse {
	req := &gate.ClientRequest{}
	resp := &gate.ClientResponse{}

	if err := proto.Unmarshal(b, req); err != nil {
		log.Error("userId : %s, error : %s", agt.userId, err.Error())
		resp.Code = base.ErrorCodeClientBadRequest

		return resp
	}

	if !strings.HasPrefix(req.Method, "/") {
		log.Error("wrong message to handle, msg: %v, userId: %s", req, agt.userId)

		return nil
	}

	resp.Id = req.Id
	resp.Method = req.Method

	// TODO: 鉴权校验、版本校验

	agt.rpcTask <- &task{
		method:  req.Method,
		content: req.Content,
		id:      req.Id,
	}

	return nil
}

// run
func (agt *Agent) run() {
	err := routine.Run(true, func() {
		defer func() {
			if !agt.closeFlag {
				// 崩溃导致的提前退出
				agt.conn.Close()
			}
		}()

		for {
			select {
			case <-agt.closeSig:
				return
			case t := <-agt.rpcTask:
				startTime := timex.Now()
				resp := &gate.ClientResponse{
					Id:     t.id,
					Method: t.method,
				}

				ctx, cancel := context.WithTimeout(context.Background(), RpcTimeout)
				methodResp, err := Invoke(GRPC, ctx, t.method, t.content)
				cancel()

				if err != nil {
					log.Error("userId : %s, req.Method : %s, isDedicatedServer : %t, error : %s", agt.userId, t.method, err.Error())
					resp.Code = base.ErrorCodeServiceBusy
				} else {
					resp.Content = methodResp
				}
				if agt.conn != nil && !agt.closeFlag {
					agt.Send(resp, agt.conn)
				}

				// 这里打印 Slow 日志，监控每个 rpc 的调用时间，提醒开发者优化（根据要求，20ms 或 40ms）
				duration := timex.Since(startTime)
				if duration > SlowThreshold {
					log.Warning("Gate slow call, userId : %s, method : %s, content: %v", agt.userId, resp.Method, t.content)
				}

				if agt.closeFlag {
					return
				}
			}
		}
	})
	if err != nil {
		agt.conn.Close()
	}
}
