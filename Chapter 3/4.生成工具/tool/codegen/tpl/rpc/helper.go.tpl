#{Insert _insert/define.i}
#{Insert _insert/header.i}
package rpc

import (
	"#{.PathBase}"
	"#{.RootPath}/internal/svc"
	"#{.RootPath}/rpc/define"
	"#{.RootPath}/rpc/grpc"
	"fmt"
)

// NewListener 创建一个服务器监听，服务端调用
func NewListener(ctx *svc.ServiceContext) define.Listener {
	// 根据配置文件的模式，选择不同的监听方式
	switch ctx.Config.RpcMode {
	case base.Grpc:
		return grpc.NewGrpcListener(ctx)
	case base.RabbitMQ:
		panic("implement me")
	case base.InterThread:
		panic("implement me")
	default:
		panic(fmt.Sprintf("unsupported rpc mode %s", ctx.Config.RpcMode))
	}
}

// New#{.PackageFirstUpper}Client 创建一个 #{.package} 客户端，其他服务创建 client 并调用任意接口
func New#{.PackageFirstUpper}Client(mode string, client *define.#{.PackageFirstUpper}) error {
	// 根据不同的模式，选择不同的连接方式
	switch mode {
	case base.Grpc:
		// TODO: 由于没有实现服务发现，所以这里需要写死地址
		// 后续可以实现基于 etcd 或 consul 的服务发现，动态填入地址
		*client = grpc.NewGrpc#{.PackageFirstUpper}("127.0.0.1:6200")
	case base.RabbitMQ:
		panic("implement me")
	case base.InterThread:
		panic("implement me")
	default:
		return fmt.Errorf("connect #{.package} service failed, unsupported runtime mode %s", mode)
	}

	return nil
}
