package router

import (
	"mygpt-back/gateway-service/internal/middleware"
	"mygpt-back/gateway-service/internal/proxy"

	"github.com/gin-gonic/gin"
)

// InitRouter 只负责注册路由
func InitUserRouter(r *gin.Engine) {
	// 认证相关 API（不需要 AuthMiddleware）
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", proxy.UserServiceProxy)
		authGroup.POST("/register", proxy.UserServiceProxy)
	}

	// 需要认证的 API
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware()) // 添加认证中间件
	{
		apiGroup.GET("/profile", proxy.UserServiceProxy) // 获取用户信息
		// apiGroup.POST("/chat", proxy.ChatServiceProxy)    // 处理聊天
	}
}
