package config

import (
	"simple/pkg/util"

	"github.com/caarlos0/env/v6"
)

var Cfg config

func InitConfig() {
	err := env.Parse(&Cfg)
	util.Assert(err == nil, "读取配置报错:%v", err)
}

type config struct {
	LoggerLevel int    `env:"LoggerLevel,required"`
	LoggerType  int    `env:"LoggerType,required"`
	ServerId    int    `env:"ServerId,required"`
	SyncFileDir string `env:"SyncFileDir,required"`
	DB          string `env:"DB,required"`
}
