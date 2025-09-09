package pluginsdk

import "context"

// Logger 插件日志接口
type Logger interface {
	// Debug 输出调试日志
	Debug(format string, args ...any)

	// Info 输出信息日志
	Info(format string, args ...any)

	// Warn 输出警告日志
	Warn(format string, args ...any)

	// Error 输出错误日志
	Error(format string, args ...any)
}

// Plugin 插件接口
type Plugin interface {
	// ID 返回插件的唯一标识
	ID() string

	// Name 返回插件名称
	Name() string

	// Version 返回插件版本
	Version() string

	// Desc 返回插件描述
	Desc() string

	// DefaultSettings 返回插件默认设置
	DefaultSettings() map[string]any

	// SetLogger 设置插件日志器
	SetLogger(logger Logger)

	// Process 处理输入数据，返回标准化输出
	Process(ctx context.Context, input map[string]any, settings map[string]any) (*Output, error)
}

// Output 插件处理输出结构
type Output struct {
	// 标题
	Title string `json:"title"`

	// 内容
	Content string `json:"content"`

	// 图片URL
	Image string `json:"image"`

	// 跳转链接
	URL string `json:"url"`

	// 多个目标
	Targets []string `json:"targets"`

	// 元数据信息
	Meta *MetaData `json:"meta"`
}

// MetaData 元数据结构
type MetaData struct {
	// 原始请求数据
	Req map[string]any `json:"req"`

	// 插件ID
	PluginID string `json:"pluginId"`

	// 处理时间戳
	ProcessedAt string `json:"processedAt"`

	// 额外元数据
	Extra map[string]any `json:"extra,omitempty"`
}

// NewFunc 插件构造函数类型
type NewFunc func() (Plugin, error)

// UIConfig UI配置结构
type UIConfig any
