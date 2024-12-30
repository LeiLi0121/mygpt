package db

import (
	"mygpt-back/pkg/config"
	"mygpt-back/pkg/logger"
	"testing"
)

func TestInitRedis(t *testing.T) {
	// 初始化 Logger
	logger.InitLogger(true)
	defer logger.CloseLogger()

	// 加载配置
	conf := config.Load("config.yaml")

	// 测试数据库初始化
	_, err := InitRedis(&conf.Redis)
	if err != nil {
		t.Fatalf("初始化 Redis 失败: %v", err)
	}
}
