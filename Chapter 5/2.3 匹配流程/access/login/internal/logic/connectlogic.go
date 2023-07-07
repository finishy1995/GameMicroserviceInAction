package logic

import (
	"ProjectX/access/login/consts"
	"ProjectX/access/login/internal/platform"
	"ProjectX/access/login/internal/types"
	"ProjectX/base"
	"ProjectX/library/log"
	"context"
)

type ConnectLogic struct {
	ctx context.Context
	log.Logger
}

func NewConnectLogic(ctx context.Context) *ConnectLogic {
	return &ConnectLogic{
		ctx:    ctx,
		Logger: log.WithContext(ctx),
	}
}

func (l *ConnectLogic) ConnectLogic(in *types.ConnectLogicRequest) (*types.ConnectLogicResponse, error) {
	if in == nil || in.Token == "" || in.Platform == "" {
		return nil, consts.ErrInvalidArgument
	}

	id, err := platform.Verify(in.Platform, in.Token)
	if err != nil {
		return &types.ConnectLogicResponse{
			Code: base.ErrorCodeInvalidPlatformOrToken,
		}, nil
	}

	// TODO: 访问 account 拿账号信息

	return &types.ConnectLogicResponse{
		Code:         base.ErrorCodeOK,
		GateEndpoint: "127.0.0.1:6200", // TODO: get gate endpoint from redis
		GateToken:    id,               // TODO: encrypt id for security
	}, nil
}
