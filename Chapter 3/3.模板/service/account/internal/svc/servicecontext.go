package svc

import "ProjectX/service/account/internal/config"

type ServiceContext struct {
	Config config.Config
}

// NewServiceContext 服务上下文，在这里连接数据库、设置 buffer、创建协程做事（例如匹配服务创建一个匹配池等）
func NewServiceContext(c config.Config) *ServiceContext {
	// TODO: 连接数据库，创建协程，设置缓冲等
	return &ServiceContext{
		Config: c,
	}
}
