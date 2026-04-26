package config

// config.go 包含了配置相关的函数，用于加载和保存用户的服务器地址和用户名

import (
	"encoding/json"
	"os"
	"strings"
)

// 默认配置文件名
const configFileName = "config.json"

// 默认公网ip
const DefaultServerAddress = "127.0.0.1:9000"

// 配置文件结构体
type Config struct {
	ServerAddress string `json:"server_address"`
	Username      string `json:"username"`
	Token         string `json:"token"`
}

// 加载配置文件
func LoadConfig() (Config, bool) {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return Config{}, false
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, false
	}

	return config, true
}

// 保存配置文件
func SaveConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFileName, data, 0644)
}

// 标准化配置,去除多余的空格,并设置默认值
func NormalizeConfig(config Config) Config {
	config.ServerAddress = strings.TrimSpace(config.ServerAddress)
	config.Username = strings.TrimSpace(config.Username)
	config.Token = strings.TrimSpace(config.Token)

	if config.ServerAddress == "" {
		config.ServerAddress = DefaultServerAddress
	}

	return config
}

// 检查配置是否ok
func IsClientConfigReady(config Config) bool {
	return strings.TrimSpace(config.ServerAddress) != "" &&
		strings.TrimSpace(config.Token) != ""
}
