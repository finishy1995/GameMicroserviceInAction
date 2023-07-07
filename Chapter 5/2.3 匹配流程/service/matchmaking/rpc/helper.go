// Code generated by CodeGenerator. DO NOT EDIT!
//
// Source: matchmaking.proto
// Time: 2023-07-05 10:20:01

package rpc

import (
	"ProjectX/base"
	"ProjectX/service/matchmaking/internal/svc"
	"ProjectX/service/matchmaking/rpc/define"
	"ProjectX/service/matchmaking/rpc/grpc"
	"fmt"
)

// NewListener 创建一个服务器监听，服务端调用
func NewListener(ctx *svc.ServiceContext) define.Listener {
	// 根据配置文件的模式，选择不同的监听方式
	switch ctx.Config.RpcMode {
	case base.Grpc:
		return grpc.NewGrpcListener(ctx)
	case base.RabbitMQ:
		panic("implement me")
	case base.InterThread:
		panic("implement me")
	default:
		panic(fmt.Sprintf("unsupported rpc mode %s", ctx.Config.RpcMode))
	}
}

// NewMatchmakingClient 创建一个 matchmaking 客户端，其他服务创建 client 并调用任意接口
func NewMatchmakingClient(mode string, client *define.Matchmaking) error {
	// 根据不同的模式，选择不同的连接方式
	switch mode {
	case base.Grpc:
		// TODO: 由于没有实现服务发现，所以这里需要写死地址
		// 后续可以实现基于 etcd 或 consul 的服务发现，动态填入地址
		*client = grpc.NewGrpcMatchmaking()
	case base.RabbitMQ:
		panic("implement me")
	case base.InterThread:
		panic("implement me")
	default:
		return fmt.Errorf("connect matchmaking service failed, unsupported runtime mode %s", mode)
	}

	return nil
}
