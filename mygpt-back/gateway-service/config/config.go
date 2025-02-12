package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    int    `mapstructure:"port"`
		Timeout string `mapstructure:"timeout"`
	} `mapstructure:"server"`

	Routes []Route `mapstructure:"routes"`

	JWT struct {
		Secret string `mapstructure:"secret"`
		Expire string `mapstructure:"expire"`
	} `mapstructure:"jwt"`

	RateLimit struct {
		RequestsPerSecond int `mapstructure:"requests_per_second"`
	} `mapstructure:"rate_limit"`
}

type Route struct {
	Path        string `mapstructure:"path"`
	Service     string `mapstructure:"service"`
	StripPrefix bool   `mapstructure:"strip_prefix"`
	Target      string `mapstructure:"target"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config_local") // 配置文件名
	viper.SetConfigType("yaml")         // 配置文件类型
	viper.AddConfigPath(configPath)     // 配置文件路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
