package proxy

import (
	"io"
	"log"
	"net/http"

	"mygpt-back/gateway-service/pkg/config"

	"github.com/gin-gonic/gin"
)

// ProxyRequest 代理请求到目标服务
func ProxyRequest(c *gin.Context, targetURL string) {
	// 创建新的 HTTP 请求
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 复制请求头
	for k, v := range c.Request.Header {
		req.Header[k] = v
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to reach target service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	// 复制响应
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// UserServiceProxy 代理请求到 `user-service`
func UserServiceProxy(c *gin.Context) {
	targetURL := config.GetConfig().Services.UserService + c.Request.URL.Path
	ProxyRequest(c, targetURL)
}

// // ChatServiceProxy 代理请求到 `chat-service`
// func ChatServiceProxy(c *gin.Context) {
// 	targetURL := config.GetConfig().Services.ChatService + c.Request.URL.Path
// 	ProxyRequest(c, targetURL)
// }
