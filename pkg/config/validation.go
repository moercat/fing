package config

import (
	"errors"
	"log"
	"strings"
)

// ValidateConfig 验证配置的有效性
func ValidateConfig() error {
	// 检查密钥是否为默认值
	if Config.Secret == "asdasggas" || strings.TrimSpace(Config.Secret) == "" {
		log.Fatal("警告: 请修改配置中的 secret 字段，不要使用默认值!")
		return errors.New("无效的密钥配置")
	}

	// 检查必要字段
	if Config.Port == "" {
		return errors.New("端口配置不能为空")
	}

	if Config.DataSource.Main == "" {
		return errors.New("数据库连接字符串不能为空")
	}

	if Config.Redis.Addr == "" {
		return errors.New("Redis地址不能为空")
	}

	return nil
}
