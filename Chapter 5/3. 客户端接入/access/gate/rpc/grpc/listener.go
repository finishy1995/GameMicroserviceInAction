// Code generated by CodeGenerator. DO NOT EDIT!
//
// Source: gate.proto
// Time: 2023-06-20 10:56:47

package grpc

import (
	"ProjectX/access/gate/internal/svc"
	"ProjectX/access/gate/rpc/define"

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
