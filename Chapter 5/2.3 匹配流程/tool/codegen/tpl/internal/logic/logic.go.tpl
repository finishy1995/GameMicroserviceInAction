#{Loop #{.service.*Length} index=.ServiceIndex}
#{Define .ServiceInstance = .service.#{.ServiceIndex}}
#{  Loop #{#{.ServiceInstance}.method.*Length} index=.ServiceMethodIndex  }
#{  Define .MethodInstance = #{.ServiceInstance}.method.#{.ServiceMethodIndex}  }
#{  StartFile  }
#{  Define file.overwrite = false  }
#{  Insert _insert/header.i  }
#{  Define .LogicName = #{#{.MethodInstance}.name}  }
#{  Define .LogicNameFirstUpper = #{Upper #{.LogicName} 1}  }
#{  Define .LogicStructName = #{.LogicNameFirstUpper}Logic  }
#{  Define file.name = internal/logic/#{Lower #{.LogicName}}logic.go  }
#{  Define .RootPath = #{.PathSuffix}#{.package} }
package logic

import (
    "#{.ProjectName}/library/log"
    "#{.RootPath}/internal/svc"
    "#{.RootPath}/pb/#{.package}"
    "context"
)

type #{.LogicStructName} struct {
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	log.Logger
}

func New#{.LogicStructName}(ctx context.Context, svcCtx *svc.ServiceContext) *#{.LogicStructName} {
	return &#{.LogicStructName}{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: log.WithContext(ctx),
	}
}

#{  Define .MethodInputType = #{.package}.#{Upper #{#{.MethodInstance}.input_type_short} 1}  }
#{  Define .MethodOutputType = #{.package}.#{Upper #{#{.MethodInstance}.output_type_short} 1}  }
func (l *#{.LogicStructName}) #{.LogicNameFirstUpper}(in *#{.MethodInputType}) (*#{.MethodOutputType}, error) {
	// TODO: logic write here

	return &#{.MethodOutputType}{}, nil
}

#{  EndFile  }
#{  EndLoop  }
#{EndLoop}