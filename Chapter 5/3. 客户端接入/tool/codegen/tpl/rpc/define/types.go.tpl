#{Insert _insert/define.i}
#{Insert _insert/header.i}
#{Insert _insert/method.i}
package define

import (
	"#{.RootPath}/pb/#{.package}"

	"context"
	"google.golang.org/grpc"
)

type (
#{Loop #{.message.*Length} index=.MessageIndex}
#{Define .MessageName = #{Upper #{.message.#{.MessageIndex}.name} 1}}
	#{.MessageName} = #{.package}.#{.MessageName}
#{EndLoop}

	#{.PackageFirstUpper} interface {
#{Loop #{.MethodLength} index=.MethodIndex}
#{Define .MethodInstance = #{.Method.#{.MethodIndex}}}
		#{Upper #{#{.MethodInstance}.name} 1}(ctx context.Context, in *#{#{.MethodInstance}.input_type_short}) (*#{#{.MethodInstance}.output_type_short}, error)
#{EndLoop}
	}

	Listener interface { // 当前服务的监听
		// Start 开启监听
		Start()
		// Stop 关闭监听
		Stop()
		// RegisterService 向 grpc 注册监听，protobuf 生成的 rpc 需要
		RegisterService(desc *grpc.ServiceDesc, impl interface{})
	}
)
