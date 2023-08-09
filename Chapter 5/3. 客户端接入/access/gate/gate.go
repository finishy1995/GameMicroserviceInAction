package main

import (
	"ProjectX/access/gate/consts"
	"ProjectX/access/gate/internal/agent"
	"ProjectX/access/gate/internal/config"
	"ProjectX/access/gate/internal/server"
	"ProjectX/access/gate/internal/svc"
	"ProjectX/access/gate/pb/gate"
	"ProjectX/access/gate/rpc"
	"ProjectX/base"
	"ProjectX/library/network"
	"ProjectX/library/network/core"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", base.GetConfigFilePathByService(consts.SvcName), "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := rpc.NewListener(ctx)
	gate.RegisterGateServer(s, server.NewGateServer(ctx))
	defer s.Stop()

	defer network.DestroyAll()
	_, err := network.Listen(network.TcpGNet, "0.0.0.0:6100", agent.NewAgent, core.WithMaxConnNum(3000))
	if err != nil {
		panic(err)
	}
	_, err = network.Listen(network.WebsocketNet, "0.0.0.0:6101", agent.NewAgent, core.WithMaxConnNum(3000))
	if err != nil {
		panic(err)
	}

	s.Start()
}
