// Code generated by CodeGenerator. Not generate if exist
//
// Source: game.proto
// Time: 2023-07-25 08:25:21

package svc

import (
	"ProjectX/access/game/internal/config"
	"ProjectX/access/game/internal/game"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	game.AddGames(c.Spec.GameSessionNum) // 启动时创建8个游戏，每个游戏有一个协程，均处于空闲状态

	// TODO: 连接数据库，创建协程，设置缓冲等
	return &ServiceContext{
		Config: c,
	}
}
