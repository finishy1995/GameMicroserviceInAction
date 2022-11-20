package define

import (
	"ProjectX/service/account/pb/account"
	"context"
	"google.golang.org/grpc"
)

type (
	AccountInfo                = account.AccountInfo
	GetOrCreateAccountRequest  = account.GetOrCreateAccountRequest
	GetOrCreateAccountResponse = account.GetOrCreateAccountResponse

	Account interface { // 给其他服务调用的方法
		// GetOrCreateAccount 获取或新账号创建
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
