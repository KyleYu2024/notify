package pluginsdk

import "context"

// Plugin 插件接口
type Plugin interface {
	// ID 返回插件的唯一标识
	ID() string

	// DefaultSettings 返回插件默认设置
	DefaultSettings() map[string]any

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
	// 是否需要通知
	IsNotify bool `json:"isNotify"`
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

type UIConfig struct {
	Component string     `json:"component"`
	Text      string     `json:"text"`
	Html      string     `json:"html"`
	Content   []UIConfig `json:"content"`
	// map[string]UIConfig || UIConfig
	Slots any `json:"slots"`
	Props any `json:"props"`
}
