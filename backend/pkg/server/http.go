package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jianxcao/notify/backend/pkg/app"
	"github.com/jianxcao/notify/backend/pkg/config"
	"github.com/jianxcao/notify/backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	app            *app.NotificationApp
	config         *config.Config
	router         *gin.Engine
	server         *http.Server
	authMiddleware *AuthMiddleware
	configManager  *config.ConfigManager
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(notificationApp *app.NotificationApp, addr string, configManager *config.ConfigManager, configFile string) *HTTPServer {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 将当前配置设置到配置管理器中
	configManager.Load() // 重新加载以确保一致性

	server := &HTTPServer{
		app:            notificationApp,
		router:         gin.New(),
		config:         configManager.GetConfig(),
		authMiddleware: NewAuthMiddleware(configManager.GetConfig()),
		configManager:  configManager,
	}

	// 添加中间件
	server.router.Use(gin.Logger())
	server.router.Use(gin.Recovery())
	server.router.Use(server.corsMiddleware())

	// 设置路由
	server.setupRoutes()

	server.server = &http.Server{
		Addr:    addr,
		Handler: server.router,
	}

	// 记录支持的应用端点
	server.logSupportedApps()

	return server
}

// corsMiddleware CORS中间件
func (s *HTTPServer) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// setupRoutes 设置所有路由
func (s *HTTPServer) setupRoutes() {
	// 设置静态文件服务 (前端资源)
	s.setupStaticRoutes()

	api := s.router.Group("/api/v1")

	// 设置健康检查路由 (定义在 health_routes.go)
	s.setupHealthRoutes(api)

	// 设置通知路由 (定义在 notify_routes.go)
	s.setupNotifyRoutes(api)

	// 设置管理路由 (定义在 admin_routes.go)
	s.setupAdminRoutes(api)

	// 设置日志流路由 (定义在 log_routes.go)
	s.setupLogRoutes(api)
}

// setupStaticRoutes 设置静态文件路由 (前端资源)
func (s *HTTPServer) setupStaticRoutes() {
	// 静态文件目录
	staticDir := config.EnvCfg.STATIC_DIR

	// 检查静态文件目录是否存在
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		logger.Warn(fmt.Sprintf("静态文件目录不存在，跳过前端资源服务: %s", staticDir))
		return
	}

	logger.Info(fmt.Sprintf("启用前端静态文件服务, dir=%s", staticDir))

	// 根路径重定向到index.html
	s.router.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, "index.html"))
	})
	s.router.HEAD("/", func(c *gin.Context) {
		// HEAD请求只返回头部，不需要文件内容
		c.Status(http.StatusOK)
	})

	// 处理所有非API路径的静态文件服务
	s.router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是API请求，返回404
		if path == "/api" || strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		fs := gin.Dir(staticDir, false)
		fileServerStatic := http.StripPrefix("/", http.FileServer(fs))
		file, err := fs.Open(c.Request.URL.Path)
		if err != nil {
			c.File(filepath.Join(staticDir, "index.html"))
			return
		} else {
			fileServerStatic.ServeHTTP(c.Writer, c.Request)
		}
		defer file.Close()
	})
}

// logSupportedApps 记录支持的应用端点
func (s *HTTPServer) logSupportedApps() {
	for appName, appConfig := range s.config.NotificationApps {
		if !appConfig.Enabled {
			continue
		}

		authStatus := "无需认证"
		if appConfig.Auth != nil && appConfig.Auth.Enabled {
			authStatus = "需要Token认证"
		}
		logger.Info(fmt.Sprintf("支持应用端点: endpoint=/api/v1/notify/%s, app_name=%s, auth_status=%s", appName, appConfig.Name, authStatus))
	}
}

// Start 启动服务器
func (s *HTTPServer) Start() error {
	logger.Info(fmt.Sprintf("启动HTTP服务器, addr=%s", s.server.Addr))
	return s.server.ListenAndServe()
}

// Stop 停止服务器
func (s *HTTPServer) Stop(ctx context.Context) error {
	logger.Info("停止HTTP服务器")
	return s.server.Shutdown(ctx)
}
