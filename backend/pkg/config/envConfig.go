package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	VERSION         string `default:"v0.0.9"`
	NOTIFY_USERNAME string
	NOTIFY_PASSWORD string
	LOG_LEVEL       string `default:"info"`
	LOG_FORMAT      string `default:"text"`
	CONFIG_FILE     string `default:"config/config.yaml"`
	PORT            string `default:":8088"`
	STATIC_DIR      string `default:"/app/static"`
}

func NewEnvConfig() *EnvConfig {
	cfg := EnvConfig{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		fmt.Println("配置加载错误")
	}
	return &cfg
}

var EnvCfg = NewEnvConfig()
