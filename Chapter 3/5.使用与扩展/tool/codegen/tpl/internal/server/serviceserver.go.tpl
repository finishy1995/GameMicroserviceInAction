#{Insert _insert/define.i}
#{Insert _insert/header.i}
#{Define file.name = internal/server/#{.package}server.go}
#{Define .ServerStructName = #{.PackageFirstUpper}Server}
package server

import (
	"context"

    "#{.RootPath}/internal/logic"
	"#{.RootPath}/internal/svc"
	"#{.RootPath}/pb/#{.package}"
)

type #{.ServerStructName} struct {
	svcCtx *svc.ServiceContext
	#{.package}.Unimplemented#{.ServerStructName}
}

func New#{.ServerStructName}(svcCtx *svc.ServiceContext) *#{.ServerStructName} {
	return &#{.ServerStructName}{
		svcCtx: svcCtx,
	}
}

#{Loop #{.service.*Length} index=.ServiceIndex}
#{Define .ServiceInstance = .service.#{.ServiceIndex}}
#{  Loop #{#{.ServiceInstance}.method.*Length} index=.ServiceMethodIndex  }
#{  Define .MethodInstance = #{.ServiceInstance}.method.#{.ServiceMethodIndex}  }
#{  Define .LogicName = #{#{.MethodInstance}.name}  }
#{  Define .LogicNameFirstUpper = #{Upper #{.LogicName} 1}  }
#{  Define .LogicStructName = #{.LogicNameFirstUpper}Logic  }
#{  Define .MethodInputType = #{.package}.#{Upper #{#{.MethodInstance}.input_type_short} 1}  }
#{  Define .MethodOutputType = #{.package}.#{Upper #{#{.MethodInstance}.output_type_short} 1}  }
func (s *#{.ServerStructName}) #{.LogicNameFirstUpper}(ctx context.Context, in *#{.MethodInputType}) (*#{.MethodOutputType}, error) {
	l := logic.New#{.LogicStructName}(ctx, s.svcCtx)
	return l.#{.LogicNameFirstUpper}(in)
}

#{  EndLoop  }
#{EndLoop}