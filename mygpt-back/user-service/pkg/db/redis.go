package db

import (
	"context"
	"fmt"
	"mygpt-back/user-service/pkg/config"
	"mygpt-back/user-service/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// InitRedis 初始化 Redis 客户端并返回实例
func InitRedis(conf *config.RedisConfig) (*redis.Client, error) {
	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password, // 如果没有密码，设置为 ""
		DB:       conf.DB,       // 默认使用的数据库
		//PoolSize: conf.PoolSize, // 连接池大小
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Log.Error("Redis 连接测试失败", zap.Error(err))
		return nil, err
	}

	// 记录成功日志
	logger.Log.Info("Redis 连接成功", zap.String("host", conf.Host), zap.Int("port", conf.Port))
	return client, nil
}

// // CloseRedis 优雅关闭 Redis 客户端
// func CloseRedis() {
// 	if RedisClient != nil {
// 		_ = RedisClient.Close()
// 		logger.Log.Info("Redis 客户端已关闭")
// 	}
// }
