package redis

import (
	"time"
)

// ClientConfig  配置信息
type ClientConfig struct {
	Host           string
	Type           string        `json:",default=node,options=node|cluster"`
	MinIdle        int           `json:",default=1"`
	MaxActive      int           `json:",default=10"`
	IdleTimeout    time.Duration `json:",default=10s"`
	Verbose        bool          `json:",default=false,options=true|false"`
	Pass           string        `json:",optional"`
	ConnectTimeout time.Duration `json:",default=100ms"`
	ReadTimeout    time.Duration `json:",default=100ms"`
	WriteTimeout   time.Duration `json:",default=100ms"`
}
