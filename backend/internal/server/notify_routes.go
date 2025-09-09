package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jianxcao/notify/internal/config"
	"github.com/jianxcao/notify/internal/logger"

	"github.com/gin-gonic/gin"
)

// setupNotifyRoutes 设置通知相关路由
func (s *HTTPServer) setupNotifyRoutes(api *gin.RouterGroup) {
	// 通知接口（使用应用Token认证）
	notify := api.Group("/notify")
	{
		notify.POST("/:appid", s.appAuthMiddleware(), s.handleSendNotification)
		notify.GET("/:appid", s.appAuthMiddleware(), s.handleSendNotificationByQuery)
		notify.PUT("/:appid", s.appAuthMiddleware(), s.handleSendNotification)
	}
}

// appAuthMiddleware 应用认证中间件（基于路径参数）
func (s *HTTPServer) appAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.Param("appid")

		// 根据 appID 查找应用配置
		var appConfig *config.NotificationApp
		var mapKey string
		found := false

		for key, app := range s.config.NotificationApps {
			if app.AppID == appID {
				appConfig = &app
				mapKey = key
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusNotFound, NewErrorRes(APP_NOT_FOUND, fmt.Sprintf("应用 %s 不存在", appID)))
			c.Abort()
			return
		}

		// 检查应用是否启用
		if !appConfig.Enabled {
			c.JSON(http.StatusForbidden, NewErrorRes(APP_DISABLED, fmt.Sprintf("应用 %s 未启用", appID)))
			c.Abort()
			return
		}

		// 检查认证（如果auth字段存在且启用）
		if appConfig.Auth != nil && appConfig.Auth.Enabled {
			if !s.authMiddleware.validateAppToken(c.Request, appConfig.Auth.Token) {
				c.JSON(http.StatusForbidden, NewErrorRes(AUTH_ERROR, "认证失败"))
				c.Abort()
				return
			}
		}

		// 将应用信息存储到上下文
		c.Set("appID", appID)
		c.Set("appConfig", *appConfig)
		c.Set("mapKey", mapKey) // 存储 map key 用于兼容现有逻辑
		c.Next()
	}
}

// ===== 通知接口处理函数 =====

// NotificationSendResponseData 通知发送响应数据结构体（驼峰命名）
type NotificationSendResponseData struct {
	AppName string `json:"appName"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Level   string `json:"level"`
	Image   string `json:"image"`
	URL     string `json:"url"`
	Method  string `json:"method"`
}

// handleSendNotification 发送通知 (POST /notify/:appname) - 从request body获取JSON数据
func (s *HTTPServer) handleSendNotification(c *gin.Context) {
	// appID := c.GetString("appID")
	appConfig := c.MustGet("appConfig").(config.NotificationApp)
	body, _ := io.ReadAll(c.Request.Body)
	contentType := c.Request.Header.Get("Content-Type")
	logger.Debug("请求体", "body", string(body))
	logger.Debug("请求头", "contentType", contentType)
	var rawData map[string]interface{} = make(map[string]interface{})
	if strings.Contains(strings.ToLower(contentType), "application/x-www-form-urlencoded") {
		formData, err := url.ParseQuery(string(body))
		if err != nil {
			logger.Error("解析请求失败", "error", err)
			c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
			return
		}
		for key, values := range formData {
			if len(values) > 0 {
				rawData[key] = values[0] // 取第一个值
			}
		}
	} else if strings.Contains(strings.ToLower(contentType), "multipart/form-data") {
		// 解析 multipart/form-data
		err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
		if err != nil {
			logger.Error("解析 multipart/form-data 失败", "error", err)
			c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析 multipart/form-data 失败"))
			return
		}

		// 获取表单字段
		if c.Request.MultipartForm != nil && c.Request.MultipartForm.Value != nil {
			for key, values := range c.Request.MultipartForm.Value {
				if len(values) > 0 {
					rawData[key] = values[0] // 取第一个值
				}
			}
		}
	} else {
		// 从JSON body获取原始数据
		err := json.Unmarshal(body, &rawData)
		if err != nil {
			logger.Error("解析请求失败", "error", err)
			c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
			return
		}
	}

	logger.Debug("发送通知原始参数", "data", rawData)

	// 发送通知
	if err := s.app.Send(c.Request.Context(), appConfig, &rawData); err != nil {
		logger.Error("发送通知失败", "error", err)
		c.JSON(http.StatusOK, NewErrorRes(NOTIFICATION_SEND_FAILED, err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, NewSuccessRes(map[string]interface{}{
		"appName": appConfig.Name,
		"method":  "POST",
	}))
}

// handleSendNotificationByQuery 发送通知 (GET /notify/:appname) - 从query参数获取
func (s *HTTPServer) handleSendNotificationByQuery(c *gin.Context) {
	// appID := c.GetString("appID")
	appConfig := c.MustGet("appConfig").(config.NotificationApp)
	// 从query参数获取原始数据
	rawData := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			rawData[key] = values[0] // 取第一个值
		}
	}
	logger.Debug("发送通知原始参数", "data", rawData)
	// 发送通知
	if err := s.app.Send(c.Request.Context(), appConfig, &rawData); err != nil {
		logger.Error("发送通知失败", "error", err)
		c.JSON(http.StatusOK, NewErrorRes(NOTIFICATION_SEND_FAILED, err.Error()))
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, NewSuccessRes(map[string]interface{}{
		"appName": appConfig.Name,
		"method":  "GET",
	}))
}
