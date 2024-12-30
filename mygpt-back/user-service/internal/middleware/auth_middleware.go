package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

// JWTClaims 定义 JWT 的载荷结构
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware 验证用户 Token 的中间件
func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization Header"})
			c.Abort()
			return
		}

		// 提取 Token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Authorization 格式"})
			c.Abort()
			return
		}

		// 验证 Token
		userID, err := validateToken(token, redisClient)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将 userID 存入上下文
		c.Set("userID", userID)

		// 执行下一个 Handler
		c.Next()
	}
}

// validateToken 验证 Token 的有效性
func validateToken(token string, redisClient *redis.Client) (int, error) {
	// 定义一个 JWTClaims 对象，用于解析 Token
	claims := &JWTClaims{}

	// 使用密钥解析 Token
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil // 替换为实际的密钥
	})
	if err != nil || !parsedToken.Valid {
		return 0, errors.New("无效的 Token")
	}

	// 检查 Token 是否已过期
	if time.Until(claims.ExpiresAt.Time) <= 0 {
		return 0, errors.New("token 已过期")
	}

	// 检查 Redis 中是否存在该 Token
	exists, err := redisClient.Exists(context.Background(), token).Result()
	if err != nil || exists == 0 {
		return 0, errors.New("token 不存在或已过期")
	}

	// 返回解析出的 userID
	return claims.UserID, nil
}
