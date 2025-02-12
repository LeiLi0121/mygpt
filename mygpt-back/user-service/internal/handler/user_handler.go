package handler

import (
	"bytes"
	"fmt"
	"io"
	"mygpt-back/user-service/internal/model"
	"mygpt-back/user-service/internal/service"
	"net/http"

	"mygpt-back/user-service/pkg/logger"
	"mygpt-back/user-service/pkg/security"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Validate(c *gin.Context) {
	// 从 gin.Context 获取 user_id (由 AuthMiddleware 设置)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}
	userIDStr := fmt.Sprintf("%v", userID) // 安全转换为字符串

	encrypted, err := security.EncryptAES(userIDStr) // 加密 user_id
	if err != nil {
		fmt.Println("加密失败:", err)
	} else {
		fmt.Println("加密后的数据:", encrypted)
	}
	// 返回成功响应，表示 Token 有效，并附带 user_id
	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"user_id": encrypted,
	})
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var user model.User
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	fmt.Println("body json is :", body)

	// 解析 JSON 请求
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		logger.Log.Error(err.Error())
		return
	}

	// 调用服务层逻辑
	if err := h.service.RegisterUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 解析 JSON 请求
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 调用服务层逻辑
	token, err := h.service.LoginUser(loginData.Username, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetProfile 获取用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取 userID（需要 Middleware 提供）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 调用服务层逻辑
	user, err := h.service.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
