package server

import (
	"context"

	"ProjectX/service/account/internal/logic"
	"ProjectX/service/account/internal/svc"
	"ProjectX/service/account/pb/account"
)

type AccountServer struct {
	svcCtx *svc.ServiceContext
	account.UnimplementedAccountServer
}

// NewAccountServer 实现 protobuf grpc 接口，串联 rpc 与 logic
func NewAccountServer(svcCtx *svc.ServiceContext) *AccountServer {
	return &AccountServer{
		svcCtx: svcCtx,
	}
}

// GetOrCreateAccount 获取或新账号创建
func (s *AccountServer) GetOrCreateAccount(ctx context.Context, in *account.GetOrCreateAccountRequest) (*account.GetOrCreateAccountResponse, error) {
	l := logic.NewGetOrCreateAccountLogic(ctx, s.svcCtx)
	return l.GetOrCreateAccount(in)
}
