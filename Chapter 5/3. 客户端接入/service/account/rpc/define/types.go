// Code generated by CodeGenerator. DO NOT EDIT!
//
// Source: account.proto
// Time: 2023-05-25 17:18:07

package define

import (
	"ProjectX/service/account/pb/account"

	"context"
	"google.golang.org/grpc"
)

type (
	GetOrCreateAccountRequest = account.GetOrCreateAccountRequest
	AccountInfo = account.AccountInfo
	GetOrCreateAccountResponse = account.GetOrCreateAccountResponse

	Account interface {
		GetOrCreateAccount(ctx context.Context, in *GetOrCreateAccountRequest) (*GetOrCreateAccountResponse, error)
	}

	Listener interface { // 当前服务的监听
		// Start 开启监听
		Start()
		// Stop 关闭监听
		Stop()
		// RegisterService 向 grpc 注册监听，protobuf 生成的 rpc 需要
		RegisterService(desc *grpc.ServiceDesc, impl interface{})
	}
)
