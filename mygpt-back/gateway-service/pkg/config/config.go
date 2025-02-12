package config

import (
	"mygpt-back/gateway-service/pkg/logger"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Config 结构体，映射 config.yaml
type Config struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Timeout string `yaml:"timeout"`
	} `yaml:"server"`

	Services struct {
		UserService string `yaml:"user_service"`
	} `yaml:"services"`

	JWT struct {
		Secret string `yaml:"secret"`
		Expire string `yaml:"expire"`
	} `yaml:"jwt"`

	RateLimit struct {
		RequestsPerSecond int `yaml:"requests_per_second"`
	} `yaml:"rate_limit"`
}

// 全局配置变量
var config *Config

// LoadConfig 读取 YAML 配置文件
func LoadConfig(filepath string) *Config {
	file, err := os.Open(filepath)
	if err != nil {
		logger.Log.Fatal("打开配置文件失败: %v", zap.Error(err))
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config = &Config{}
	if err := decoder.Decode(config); err != nil {
		logger.Log.Fatal("解析配置文件失败: %v", zap.Error(err))
	}

	logger.Log.Info("配置文件加载成功", zap.String("filepath", filepath))
	return config
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if config == nil {
		logger.Log.Fatal("配置未初始化，请先调用 LoadConfig()")
	}
	return config
}
