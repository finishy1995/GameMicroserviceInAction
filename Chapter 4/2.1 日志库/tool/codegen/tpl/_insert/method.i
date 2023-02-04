#{Define .MethodLength = 0}
#{Loop #{.service.*Length} index=.ServiceIndex}
#{  Define .ServiceInstance = .service.#{.ServiceIndex}  }
#{  Loop #{#{.ServiceInstance}.method.*Length} index=.ServiceMethodIndex  }
#{    Define .MethodLength = #{Calc #{.MethodLength} + 1}    }
#{    Define .Method.#{.MethodLength} = #{.ServiceInstance}.method.#{.ServiceMethodIndex}    }
#{  EndLoop  }
#{EndLoop}