#{Insert _insert/define.i}
#{Define file.overwrite = false}
#{Insert _insert/header.i}
package svc

import "#{.RootPath}/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	// TODO: 连接数据库，创建协程，设置缓冲等
	return &ServiceContext{
		Config: c,
	}
}
