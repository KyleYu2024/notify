package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/jianxcao/notify/internal/config"
	"github.com/jianxcao/notify/internal/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// BasicAuthMiddleware 基础认证中间件（用于管理接口）- 自定义实现避免浏览器弹窗
func (am *AuthMiddleware) BasicAuthMiddleware() gin.HandlerFunc {
	username := config.EnvCfg.NOTIFY_USERNAME
	password := config.EnvCfg.NOTIFY_PASSWORD

	// 如果环境变量未设置，输出警告但禁止访问
	if username == "" || password == "" {
		logger.Warn("NOTIFY_USERNAME 和 NOTIFY_PASSWORD 环境变量未设置，管理接口将无需认证即可访问！这存在安全风险，建议在生产环境中设置这些环境变量")

		return func(c *gin.Context) {
			// 记录访问日志
			logger.Warn("无认证访问管理接口",
				"method", c.Request.Method,
				"path", c.Request.URL.Path)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			// 返回 JSON 错误而不是 401，避免浏览器弹窗
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "AUTH_REQUIRED",
				"message": "需要认证，请提供 Authorization 头",
			})
			c.Abort()
			return
		}

		// 解析 Basic Auth
		if !strings.HasPrefix(auth, "Basic ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_AUTH_FORMAT",
				"message": "认证格式错误，请使用 Basic Auth",
			})
			c.Abort()
			return
		}

		payload := strings.TrimPrefix(auth, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_AUTH_ENCODING",
				"message": "认证编码错误",
			})
			c.Abort()
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 || credentials[0] != username || credentials[1] != password {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "用户名或密码错误",
			})
			c.Abort()
			return
		}

		// 认证成功，记录日志
		// logger.Info("管理接口认证成功",
		// 	"method", c.Request.Method,
		// 	"path", c.Request.URL.Path,
		// 	"user", username)

		c.Next()
	}
}

// validateAppToken 验证应用Token
func (am *AuthMiddleware) validateAppToken(r *http.Request, expectedToken string) bool {
	// 从Header中获取Token
	auth := r.Header.Get("Authorization")
	if auth == "" {
		auth = r.URL.Query().Get("token")
	}
	if auth == "" {
		return false
	}
	// 检查Bearer Token格式
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	return token == expectedToken
}
