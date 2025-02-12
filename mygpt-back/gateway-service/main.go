package main

import (
	"log"
	"strconv"

	"mygpt-back/gateway-service/internal/router"
	"mygpt-back/gateway-service/pkg/config"
	"mygpt-back/gateway-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1️⃣ 加载配置
	logger.InitLogger(true)
	// config.LoadConfig("config/config_local.yaml")
	config.LoadConfig("config/config_docker.yaml")
	// 2️⃣ 初始化 Gin
	r := gin.Default()

	// 3️⃣ 注册路由
	router.InitUserRouter(r)

	// 4️⃣ 启动 Gateway Service
	port := config.GetConfig().Server.Port
	portString := config.GetConfig().Server.Host + ":" + strconv.Itoa(port) // 将端口号转换为字符串
	log.Printf("🚀 Gateway 服务启动，监听端口 %s", portString)
	if err := r.Run(portString); err != nil {
		log.Fatalf("❌ 服务启动失败: %v", err)
	}
}
