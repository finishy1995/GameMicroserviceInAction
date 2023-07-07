#{Insert _insert/define.i}
#{Define file.overwrite = false}
#{Insert _insert/header.i}
#{Define file.name = #{.package}.go}
package main

import (
	"#{.PathBase}"
	"#{.RootPath}/consts"
	"#{.RootPath}/internal/config"
	"#{.RootPath}/internal/server"
	"#{.RootPath}/internal/svc"
	"#{.RootPath}/pb/#{.package}"
	"#{.RootPath}/rpc"
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
	#{.package}.Register#{.PackageFirstUpper}Server(s, server.New#{.PackageFirstUpper}Server(ctx))
	defer s.Stop()
	s.Start()
}
