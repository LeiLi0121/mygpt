package db

import (
	"fmt"
	"mygpt-back/pkg/config"
	"mygpt-back/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL(conf *config.MySQLConfig) (*gorm.DB, error) {
	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)

	// 使用封装的 GORM Logger
	gormLogger := logger.NewZapGormLogger(logger.Log)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger, // 使用封装后的日志器
	})
	if err != nil {
		logger.Log.Error("无法连接到 MySQL 数据库", zap.Error(err))
		return nil, fmt.Errorf("无法连接到 MySQL 数据库: %w", err)
	}

	logger.Log.Info("MySQL 数据库连接成功", zap.String("host", conf.Host), zap.Int("port", conf.Port))
	return db, nil
}
