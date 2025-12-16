package config

import (
	"github.com/jinzhu/configor"
	"log"
)

func init() {
	err := configor.New(&configor.Config{
		Debug: true,
	}).Load(&Config, "config.yaml")
	if err != nil {
		log.Fatalf("配置加载错误: %v", err)
	}

	// 验证配置的有效性
	if err := ValidateConfig(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}
}

var Config = struct {
	Mode   string `yaml:"mode" default:"debug"`
	Secret string `yaml:"secret" default:"djaod"`
	Level  int    `yaml:"level" default:"4"`
	Port   string `yaml:"port" required:"true"`

	DataSource struct {
		Main string `yaml:"main" required:"true"`
	} `yaml:"dataSource"`

	Redis struct {
		Addr string `yaml:"addr" required:"true"`
		Pool struct {
			MaxIdle     int `yaml:"maxIdle" default:"5"`
			MaxActive   int `yaml:"maxActive" default:"10"`
			IdleTimeout int `yaml:"idleTimeout" default:"180"`
		} `yaml:"pool"`
	} `yaml:"redis"`

	Es struct {
		EsUrl      string `yaml:"esUrl" default:""`
		EsUsername string `yaml:"esUsername"`
		EsPassword string `yaml:"esPassword"`
	} `yaml:"es"`

	Email struct {
		Host     string `yaml:"host" default:""`
		Name     string `yaml:"name" default:""`
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
	} `yaml:"email"`
}{}
