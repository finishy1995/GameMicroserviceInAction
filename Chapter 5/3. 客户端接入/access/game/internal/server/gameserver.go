// Code generated by CodeGenerator. DO NOT EDIT!
//
// Source: game.proto
// Time: 2023-07-25 10:14:49

package server

import (
	"context"

    "ProjectX/access/game/internal/logic"
	"ProjectX/access/game/internal/svc"
	"ProjectX/access/game/pb/game"
)

type GameServer struct {
	svcCtx *svc.ServiceContext
	game.UnimplementedGameServer
}

func NewGameServer(svcCtx *svc.ServiceContext) *GameServer {
	return &GameServer{
		svcCtx: svcCtx,
	}
}

func (s *GameServer) SetGameEnvironment(ctx context.Context, in *game.SetGameEnvironmentRequest) (*game.SetGameEnvironmentResponse, error) {
	l := logic.NewSetGameEnvironmentLogic(ctx, s.svcCtx)
	return l.SetGameEnvironment(in)
}
