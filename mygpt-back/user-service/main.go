package main

import (
	"log"
	"mygpt-back/user-service/internal/handler"
	"mygpt-back/user-service/internal/model"
	"mygpt-back/user-service/pkg/config"
	"mygpt-back/user-service/pkg/db"
	"mygpt-back/user-service/pkg/logger"
	"strconv"

	middleware "mygpt-back/user-service/internal/middleware"
	"mygpt-back/user-service/internal/repository"
	"mygpt-back/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	//debug mode on
	logger.InitLogger(true) // 如果需要自定义日志，可以实现这个函数

	// 加载配置
	// config := config.Load("config/config_local.yaml")
	config := config.Load("config/config_docker.yaml")
	log.Println("配置文件加载成功")

	// 初始化数据库
	mysqlDB, err := db.InitMySQL(&config.MySQL)
	if err != nil {
		log.Fatalf("MySQL 初始化失败: %v", err)
	}

	// 自动迁移
	if err := mysqlDB.AutoMigrate(&model.User{}); err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
	log.Println("MySQL 初始化成功")

	// 初始化 Redis
	redisClient, err := db.InitRedis(&config.Redis)
	if err != nil {
		log.Fatalf("Redis 初始化失败: %v", err)
	}
	log.Println("Redis 初始化成功")

	// 创建服务和处理器实例
	userRepo := repository.NewUserRepository(mysqlDB)
	userService := service.NewUserService(userRepo, redisClient)
	userHandler := handler.NewUserHandler(userService)

	// 配置路由
	r := gin.Default()

	// 无需验证的路由
	authGroup := r.Group("/api/auth") // authentification api
	{
		authGroup.POST("/register", userHandler.Register) //需要做格式限制和检查
		authGroup.POST("/login", userHandler.Login)       //需要做格式限制和检查; 邮箱登录和账号登录的切换;
	}
	// 需要验证的路由（使用 Group）
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(redisClient)) // 添加中间件
	{
		// ✅ **新增: gateway-service 通过这个接口验证 Token**
		auth.POST("/validate", userHandler.Validate)

		auth.GET("/profile", userHandler.GetProfile)
	}
	// 启动服务器
	port := config.Server.Port
	if port == 0 {
		port = 8081 // 默认端口
	}
	portString := config.Server.Host + ":" + strconv.Itoa(port) // 将端口号转换为字符串
	log.Printf("服务启动中，监听端口 %s...", portString)
	if err := r.Run(portString); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
