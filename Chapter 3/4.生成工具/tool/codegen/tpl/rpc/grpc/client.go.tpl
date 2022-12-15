#{Insert _insert/define.i}
#{Insert _insert/header.i}
#{Insert _insert/method.i}
package grpc

import (
	"#{.RootPath}/pb/#{.package}"
	"#{.RootPath}/rpc/define"

	"context"
	"google.golang.org/grpc"
)

type grpc#{.PackageFirstUpper} struct {
	client *grpc.ClientConn
}

func NewGrpc#{.PackageFirstUpper}(address string) define.#{.PackageFirstUpper} {
	client, err := grpc.Dial(address)
	if err != nil {
		return nil
	}
	return &grpc#{.PackageFirstUpper}{
		client: client,
	}
}

#{Loop #{.MethodLength} index=.MethodIndex}
#{Define .MethodInstance = #{.Method.#{.MethodIndex}}}
func (r *grpc#{.PackageFirstUpper}) #{Upper #{#{.MethodInstance}.name} 1}(ctx context.Context, in *define.#{#{.MethodInstance}.input_type_short}) (*define.#{#{.MethodInstance}.output_type_short}, error) {
	client := #{.package}.New#{.PackageFirstUpper}Client(r.client)
	return client.#{Upper #{#{.MethodInstance}.name} 1}(ctx, in)
}

#{EndLoop}