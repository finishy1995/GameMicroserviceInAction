package grpc

import (
	"ProjectX/service/account/pb/account"
	"ProjectX/service/account/rpc/define"
	"context"
	"google.golang.org/grpc"
)

type grpcAccount struct {
	client *grpc.ClientConn
}

// NewGrpcAccount 创建一个 account 服务的逻辑调用端，其他服务可以直接使用
// Example:
//  1. account := NewGrpcAccount() 获取 Account 逻辑调用端
//  2. account.GetOrCreateAccount() 调用具体逻辑
//
// Notes: 以下代码仅为示例，未实现重试，服务发现，负载均衡，熔断等特性
// 可以使用 zrpc 等 grpc 上层实现，获取更完善的连接
// 其他调用方法同理，RabbitMQ 等封装出来类似，内部实现更复杂
func NewGrpcAccount(address string) define.Account {
	client, err := grpc.Dial(address)
	if err != nil {
		return nil
	}
	return &grpcAccount{
		client: client,
	}
}

// GetOrCreateAccount 获取或新账号创建
func (m *grpcAccount) GetOrCreateAccount(ctx context.Context, in *define.GetOrCreateAccountRequest) (*define.GetOrCreateAccountResponse, error) {
	client := account.NewAccountClient(m.client)
	return client.GetOrCreateAccount(ctx, in)
}
