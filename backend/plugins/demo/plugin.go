package main

import (
	"context"
	pluginsdk "demo-plugin/internal"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

var baseLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true, // ✅ 打开源码位置
}))
var logger = baseLogger.WithGroup("demo")

// DemoPlugin 演示插件
type DemoPlugin struct {
}

// ID 返回插件唯一标识
func (p *DemoPlugin) ID() string {
	return "demo"
}

// Name 返回插件名称
func (p *DemoPlugin) Name() string {
	return "Demo Plugin"
}

// Version 返回插件版本
func (p *DemoPlugin) Version() string {
	return "1.0.0"
}

// Desc 返回插件描述
func (p *DemoPlugin) Desc() string {
	return "演示插件，展示如何处理通知数据并转换为标准格式"
}

// DefaultSettings 返回插件默认设置
func (p *DemoPlugin) DefaultSettings() map[string]any {
	return map[string]any{
		"prefix":        "Demo",
		"add_timestamp": true,
		"default_image": "https://example.com/demo.png",
		"debug":         false,
	}
}

// Process 处理输入数据，返回标准化输出
func (p *DemoPlugin) Process(ctx context.Context, input map[string]any, settings map[string]any) (*pluginsdk.Output, error) {
	// 记录处理开始
	logger.Info("开始处理通知数据")

	// 从输入数据中提取信息
	title, ok := input["title"].(string)
	if !ok {
		title = "Demo"
	}
	title = fmt.Sprintf("Demo: %s", title)
	content, ok := input["content"].(string)
	if !ok {
		content = "Demo"
	}
	image, ok := input["image"].(string)
	if !ok {
		image = ""
	}
	url, ok := input["url"].(string)
	if !ok {
		url = ""
	}

	logger.Debug(fmt.Sprintf("提取数据: title='%s', content='%s', image='%s', url='%s'", title, content, image, url))

	// 处理目标列表
	targetsStr, ok := input["targets"].(string)
	targets := []string{}
	if ok && targetsStr != "" {
		targets = strings.Split(targetsStr, ",")
	}

	if len(targets) > 0 {
		logger.Debug(fmt.Sprintf("找到目标列表: %v", targets))
	}

	// 构建输出
	output := &pluginsdk.Output{
		Title:   title,
		Content: content,
		Image:   image,
		URL:     url,
		Targets: targets,
		Meta: &pluginsdk.MetaData{
			Req:         input,
			PluginID:    p.ID(),
			ProcessedAt: time.Now().Format(time.RFC3339),
			Extra: map[string]any{
				"settings": settings,
			},
		},
	}

	// 记录处理完成
	logger.Info(fmt.Sprintf("通知数据处理完成: title='%s', targets=%d个", output.Title, len(output.Targets)))

	return output, nil
}

// NewPlugin 插件构造函数（插件入口点）
func NewPlugin() (pluginsdk.Plugin, error) {
	return &DemoPlugin{}, nil
}

// main 函数是必需的，但不会被调用
func main() {
	// 这个函数不会被调用，只是为了让 Go 能够编译成插件
}
