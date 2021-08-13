package config

import (
	"simple/lib/util"

	"github.com/caarlos0/env/v6"
)

var Cfg config

func init() {
	err := env.Parse(&Cfg)
	util.Assert(err == nil, "读取配置报错:%v", err)
}

type config struct {
	LoggerLevel int `env:"LoggerLevel,required"`
	LoggerType  int `env:"LoggerType,required"`
}
