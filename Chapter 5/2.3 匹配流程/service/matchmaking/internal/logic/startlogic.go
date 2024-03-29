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

type StartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	log.Logger
}

func NewStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartLogic {
	return &StartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: log.WithContext(ctx),
	}
}

func (l *StartLogic) Start(_ *matchmaking.StartRequest) (*matchmaking.StartResponse, error) {
	resp := &matchmaking.StartResponse{
		Code: base.ErrorCodeOK,
	}

	userId := contextx.GetValueFromContext(l.ctx, base.UserId)
	if userId == "" {
		resp.Code = base.ErrorCodeInternalError
		return resp, nil
	}
	ticketId := l.svcCtx.Pool.AddUser(userId)
	if ticketId == "" {
		resp.Code = base.ErrorCodeServiceBusy
	}

	return resp, nil
}
