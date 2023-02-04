#{Define file.overwrite = false}
#{Define .ServiceConfigName = #{.package}SpecialConfig}
#{Insert _insert/header.i}
package config

type RpcServerBaseConfig struct {
	ListenOn string
	RpcMode  string `json:",default=grpc,options=grpc|rabbit|inter"`
}

type Config struct {
	RpcServerBaseConfig
	Spec #{.ServiceConfigName}
}

type #{.ServiceConfigName} struct {
}
