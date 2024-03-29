// Code generated by CodeGenerator. Not generate if exist
//
// Source: matchmaking.proto
// Time: 2023-07-05 10:20:01

package logic

import (
	"ProjectX/base"
	"ProjectX/library/contextx"
	"ProjectX/library/log"
	"ProjectX/service/matchmaking/internal/svc"
	"ProjectX/service/matchmaking/pb/matchmaking"
	"context"
)

type CancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	log.Logger
}

func NewCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLogic {
	return &CancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: log.WithContext(ctx),
	}
}

func (l *CancelLogic) Cancel(_ *matchmaking.CancelRequest) (*matchmaking.CancelResponse, error) {
	resp := &matchmaking.CancelResponse{
		Code: base.ErrorCodeOK,
	}

	userId := contextx.GetValueFromContext(l.ctx, base.UserId)
	if userId == "" {
		resp.Code = base.ErrorCodeInternalError
		return resp, nil
	}
	l.svcCtx.Pool.CancelUser(userId)

	return resp, nil
}
