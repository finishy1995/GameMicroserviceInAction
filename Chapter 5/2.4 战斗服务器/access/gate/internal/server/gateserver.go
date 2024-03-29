// Code generated by CodeGenerator. DO NOT EDIT!
//
// Source: gate.proto
// Time: 2023-06-20 10:56:47

package server

import (
	"context"

	"ProjectX/access/gate/internal/logic"
	"ProjectX/access/gate/internal/svc"
	"ProjectX/access/gate/pb/gate"
)

type GateServer struct {
	svcCtx *svc.ServiceContext
	gate.UnimplementedGateServer
}

func NewGateServer(svcCtx *svc.ServiceContext) *GateServer {
	return &GateServer{
		svcCtx: svcCtx,
	}
}

func (s *GateServer) GetUserNumber(ctx context.Context, in *gate.GetUserNumberRequest) (*gate.GetUserNumberResponse, error) {
	l := logic.NewGetUserNumberLogic(ctx, s.svcCtx)
	return l.GetUserNumber(in)
}
