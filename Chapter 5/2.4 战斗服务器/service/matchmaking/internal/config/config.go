// Code generated by CodeGenerator. Not generate if exist
//
// Source: matchmaking.proto
// Time: 2023-07-05 10:20:01

package config

type RpcServerBaseConfig struct {
	ListenOn string
	RpcMode  string `json:",default=grpc,options=grpc|rabbit|inter"`
}

type Config struct {
	RpcServerBaseConfig
	Spec matchmakingSpecialConfig
}

type matchmakingSpecialConfig struct {
}
