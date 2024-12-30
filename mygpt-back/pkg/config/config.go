package config

import (
	"mygpt-back/pkg/logger"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Config 结构体，包含 server、mysql 和 redis 配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
}

// ServerConfig 表示服务配置
type ServerConfig struct {
	Host string `yaml:"host"` // 服务绑定地址
	Port int    `yaml:"port"` // 服务端口
}

// MySQLConfig 表示 MySQL 数据库配置
type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// RedisConfig 表示 Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Load 从指定文件路径加载配置
func Load(filepath string) *Config {
	file, err := os.Open(filepath)
	if err != nil {
		// 直接记录错误，退出程序
		logger.Log.Fatal("加载配置文件失败", zap.String("filepath", filepath), zap.Error(err))
	}
	defer file.Close()

	var conf Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		// 记录解析错误，退出程序
		logger.Log.Fatal("解析配置文件失败", zap.String("filepath", filepath), zap.Error(err))
	}
	return &conf
}
