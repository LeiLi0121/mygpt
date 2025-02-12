package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"mygpt-back/gateway-service/pkg/config"
	"mygpt-back/gateway-service/pkg/logger"
	"mygpt-back/gateway-service/pkg/security"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TokenValidationResponse 结构体
type TokenValidationResponse struct {
	Valid       bool   `json:"valid"`
	EncryptedID string `json:"user_id"` // ✅ 确保字段名与 user-service 返回的一致
}

// AuthMiddleware 通过 HTTP 调用 user-service 进行 Token 验证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头部
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// 1️⃣ 调用 user-service 进行 Token 验证
		userServiceURL := config.GetConfig().Services.UserService + "/api/validate"

		// 发送 HTTP 请求（Token 作为 Header）
		req, err := http.NewRequest("POST", userServiceURL, nil)
		if err != nil {
			log.Printf("AuthMiddleware: 创建 HTTP 请求失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		req.Header.Set("Authorization", token)

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("AuthMiddleware: 无法访问 user-service: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service unavailable"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// 解析响应
		var validationResp TokenValidationResponse
		if err := json.NewDecoder(resp.Body).Decode(&validationResp); err != nil {
			logger.Log.Error("AuthMiddleware: 解析 user-service 响应失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from authentication service"})
			c.Abort()
			return
		}

		// 2️⃣ 判断 Token 是否有效
		if !validationResp.Valid {
			// logger.Log.Info("respond:",validationResp)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 3️⃣ 解密 `user_id`
		decryptedUserID, err := security.DecryptAES(validationResp.EncryptedID)
		if err != nil {
			logger.Log.Error("AuthMiddleware: 解密 user_id 失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt user ID"})
			c.Abort()
			return
		}

		// 4️⃣ 存入 Context，供后续处理
		c.Set("user_id", decryptedUserID)

		// 继续请求
		c.Next()
	}
}
