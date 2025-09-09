package server

import (
	"fmt"
	"net/http"

	"github.com/jianxcao/notify/internal/config"

	"github.com/gin-gonic/gin"
)

// setupHealthRoutes 设置健康检查路由
func (s *HTTPServer) setupHealthRoutes(api *gin.RouterGroup) {
	// 健康检查（无需认证）
	api.GET("/health", s.handleHealth)
	api.HEAD("/health", s.handleHealth)
}

// SupportedApp 支持的应用信息结构体（驼峰命名）
type SupportedApp struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Endpoint    string `json:"endpoint"`
	AuthEnabled bool   `json:"authEnabled"`
}

// AdminEndpoints 管理端点结构体（驼峰命名）
type AdminEndpoints struct {
	Apps      string `json:"apps"`
	Notifiers string `json:"notifiers"`
	Config    string `json:"config"`
}

// HealthData 健康检查数据结构体（驼峰命名）
type HealthData struct {
	Status            string         `json:"status"`
	SupportedApps     []SupportedApp `json:"supportedApps"`
	AdminEndpoints    AdminEndpoints `json:"adminEndpoints"`
	AdminAuthRequired bool           `json:"adminAuthRequired"` // 新增：admin接口是否需要认证
	Version           string         `json:"version"`
}

// handleHealth 健康检查
func (s *HTTPServer) handleHealth(c *gin.Context) {
	// 验证配置
	if err := s.app.ValidateConfig(); err != nil {
		c.JSON(http.StatusServiceUnavailable, NewErrorRes(HEALTH_CHECK_FAILED, err.Error()))
		return
	}

	// 构建健康检查响应，包含所有支持的应用端点
	supportedApps := make([]SupportedApp, 0)
	for appName, appConfig := range s.config.NotificationApps {
		if appConfig.Enabled {
			authEnabled := false
			if appConfig.Auth != nil {
				authEnabled = appConfig.Auth.Enabled
			}

			supportedApps = append(supportedApps, SupportedApp{
				Name:        appName,
				DisplayName: appConfig.Name,
				Endpoint:    fmt.Sprintf("/api/v1/notify/%s", appName),
				AuthEnabled: authEnabled,
			})
		}
	}

	// 检查admin是否需要认证：如果环境变量设置了用户名和密码，则需要认证
	adminAuthRequired := config.EnvCfg.NOTIFY_USERNAME != "" && config.EnvCfg.NOTIFY_PASSWORD != ""

	healthData := HealthData{
		Status:            "healthy",
		SupportedApps:     supportedApps,
		AdminAuthRequired: adminAuthRequired,
		AdminEndpoints: AdminEndpoints{
			Apps:      "/api/v1/admin/apps",
			Notifiers: "/api/v1/admin/notifiers",
			Config:    "/api/v1/admin/config",
		},
		Version: config.EnvCfg.VERSION,
	}

	c.JSON(http.StatusOK, NewSuccessRes(healthData))
}
