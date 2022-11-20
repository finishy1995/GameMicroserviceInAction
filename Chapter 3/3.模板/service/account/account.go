package main

import (
	"ProjectX/base"
	"ProjectX/service/account/consts"
	"ProjectX/service/account/internal/config"
	"ProjectX/service/account/internal/server"
	"ProjectX/service/account/internal/svc"
	"ProjectX/service/account/pb/account"
	"ProjectX/service/account/rpc"
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", base.GetConfigFilePathByService(consts.SvcName), "the config file")

func main() {
	flag.Parse() // 加载命令行参数

	var c config.Config
	conf.MustLoad(*configFile, &c)  // 加载配置文件
	ctx := svc.NewServiceContext(c) // 初始化服务上下文

	s := rpc.NewListener(ctx)                                      // 创建服务监听
	account.RegisterAccountServer(s, server.NewAccountServer(ctx)) // 向 grpc 注册服务监听
	defer s.Stop()                                                 // account 服务关闭前执行，关闭服务监听
	s.Start()                                                      // 开启服务监听，服务器开始运行并处理请求
}
