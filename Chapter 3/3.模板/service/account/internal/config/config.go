package config

// RpcServerBaseConfig 基本配置项，后续需要移动到 library 公共库中，共享给所有服务
type RpcServerBaseConfig struct {
	ListenOn string // 监听地址
	RpcMode  string `json:",default=grpc,options=grpc|rabbit|inter"` // rpc 模式
	// 可以把数据库、Etcd、RabbitMQ、超时时间等服务间的通用配置都写在这里
}

// Config Account 服务的配置
type Config struct {
	RpcServerBaseConfig                      // 加载通用配置项
	Spec                AccountSpecialConfig // Account 服务独特的配置项
}

// AccountSpecialConfig Account 服务独特的配置项
type AccountSpecialConfig struct {
	MaxAccountNum int32 `json:",default=-1"` // 示例：最大创建账号数量，配置默认值为 -1 不限制
}
