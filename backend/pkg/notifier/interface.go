package notifier

import "context"

// NotificationMessage 通知消息结构
type NotificationMessage struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Image     string `json:"image"` // 图片URL或路径
	URL       string `json:"url"`   // 点击跳转的URL
}

// Notifier 通知服务接口
type Notifier interface {
	// Name 返回通知服务的名称
	Name() string

	// Send 发送通知消息
	Send(ctx context.Context, message *NotificationMessage, targets []string) error

	// IsEnabled 检查服务是否启用
	IsEnabled() bool

	// Validate 验证配置是否有效
	Validate() error
}

// NotificationTarget 通知目标
type NotificationTarget struct {
	Type string `json:"type"` // user, group, channel等
	ID   string `json:"id"`   // 目标ID
	Name string `json:"name"` // 目标名称
}
