package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/jianxcao/notify/internal/app"
	"github.com/jianxcao/notify/internal/config"
	"github.com/jianxcao/notify/internal/logger"
	"github.com/jianxcao/notify/internal/server"
)

var (
	configFile = flag.String("config", "config/config.yaml", "配置文件路径")
	addr       = flag.String("addr", ":8080", "HTTP服务器监听地址")
	version    = flag.Bool("version", false, "显示版本信息")
)

const (
	appVersion = "1.0.0"
	appName    = "notify"
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		return
	}
	fmt.Println("当前Go版本:", runtime.Version())
	// 获取配置文件路径，优先使用环境变量，如果环境变量没有设置则使用命令行参数
	var actualConfigFile string
	if config.EnvCfg.CONFIG_FILE != "" {
		actualConfigFile = config.EnvCfg.CONFIG_FILE
	} else {
		actualConfigFile = *configFile
	}

	// 加载配置
	configManager := config.NewConfigManager(actualConfigFile)
	_, err := configManager.Load()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志系统
	logger.Init()

	// 创建通知应用
	notificationApp := app.NewNotificationApp(configManager)

	// 验证通知应用配置
	if err := notificationApp.ValidateConfig(); err != nil {
		logger.Fatal("通知应用配置验证失败", "error", err)
	}

	// 记录环境变量认证状态
	username := config.EnvCfg.NOTIFY_USERNAME
	password := config.EnvCfg.NOTIFY_PASSWORD
	if username != "" && password != "" {
		logger.Info("Basic认证已启用", "username", username)
	} else {
		logger.Warn("Basic认证未配置，管理接口将无需认证即可访问 - 存在安全风险！")
		logger.Warn("建议设置 NOTIFY_USERNAME 和 NOTIFY_PASSWORD 环境变量以启用认证")
	}

	// 获取监听地址，优先使用环境变量，如果环境变量没有设置则使用命令行参数
	var actualAddr string
	if config.EnvCfg.PORT != "" {
		actualAddr = ":" + strings.TrimPrefix(config.EnvCfg.PORT, ":")
	} else {
		actualAddr = *addr
	}

	// 创建HTTP服务器，传递配置文件路径
	httpServer := server.NewHTTPServer(notificationApp, actualAddr, configManager, actualConfigFile)

	// 启动服务器
	go func() {
		if err := httpServer.Start(); err != nil {
			logger.Fatal("启动HTTP服务器失败", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		logger.Fatal("强制关闭服务器", "error", err)
	}

	logger.Info("服务器已关闭")
}
