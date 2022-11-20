package grpc

import (
	"ProjectX/service/account/internal/svc"
	"ProjectX/service/account/rpc/define"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type grpcListener struct {
	s       *grpc.Server
	address string
}

func NewGrpcListener(ctx *svc.ServiceContext) define.Listener {
	s := grpc.NewServer()
	return &grpcListener{
		s:       s,
		address: ctx.Config.ListenOn,
	}
}

func (listener *grpcListener) Start() {
	netListener, err := net.Listen("tcp", listener.address)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Starting rpc server at %s...\n", listener.address)
	err = listener.s.Serve(netListener)
	if err != nil {
		panic(err)
	}
}

func (listener *grpcListener) Stop() {
	listener.s.Stop()
}

func (listener *grpcListener) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	listener.s.RegisterService(desc, impl)
}
