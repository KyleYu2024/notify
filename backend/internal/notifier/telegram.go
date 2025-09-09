package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/jianxcao/notify/internal/config"
	"github.com/jianxcao/notify/internal/logger"

	"github.com/go-resty/resty/v2"
)

// TelegramNotifier Telegram通知服务
type TelegramNotifier struct {
	config config.TelegramConfig
	client *resty.Client
}

// TelegramResponse Telegram API响应结构
type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
}

// NewTelegramNotifier 创建Telegram通知服务实例
func NewTelegramNotifier(cfg config.TelegramConfig) *TelegramNotifier {
	client := resty.New()
	client.SetTimeout(100 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(10 * time.Second)

	// 如果配置了代理，设置代理
	if cfg.Proxy != "" {
		client.SetProxy(cfg.Proxy)
	}

	return &TelegramNotifier{
		config: cfg,
		client: client,
	}
}

// Name 返回服务名称
func (t *TelegramNotifier) Name() string {
	return "telegram"
}

// IsEnabled 检查服务是否启用
func (t *TelegramNotifier) IsEnabled() bool {
	return t.config.Enabled
}

// Validate 验证配置
func (t *TelegramNotifier) Validate() error {
	if !t.config.Enabled {
		return nil
	}

	if t.config.BotToken == "" {
		return fmt.Errorf("telegram Bot Token 不能为空")
	}
	if t.config.ChatID == "" {
		return fmt.Errorf("telegram Chat ID 不能为空")
	}

	return nil
}

// Send 发送通知消息
func (t *TelegramNotifier) Send(ctx context.Context, message *NotificationMessage, targets []string) error {
	if !t.config.Enabled {
		return fmt.Errorf("Telegram通知服务未启用")
	}
	users := []string{t.config.ChatID}
	if len(targets) > 0 {
		users = targets
	}
	for _, user := range users {
		// 使用message.Targets作为@提及的用户名列表
		if message.Image != "" {
			err := t.sendPhotoMessage(ctx, user, message)
			if err != nil {
				return fmt.Errorf("发送图片消息失败: %w", err)
			}
		} else {
			err := t.sendTextMessage(ctx, user, message)
			if err != nil {
				return fmt.Errorf("发送文本消息失败: %w", err)
			}
		}
	}

	return nil
}

// sendTextMessage 发送文本消息
func (t *TelegramNotifier) sendTextMessage(ctx context.Context, chatID string, message *NotificationMessage) error {
	content := fmt.Sprintf("*%s*\n\n%s", message.Title, message.Content)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.config.BotToken)
	// TODO: 消息超过4096字符的处理
	requestBody := map[string]interface{}{
		"chat_id":    chatID,
		"text":       content,
		"parse_mode": "Markdown",
	}

	// 如果有URL，添加inline keyboard按钮
	if message.URL != "" {
		requestBody["reply_markup"] = map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{
						"text": "🔗 查看详情",
						"url":  message.URL,
					},
				},
			},
		}
	}

	return t.sendRequest(ctx, apiURL, requestBody)
}

// sendPhotoMessage 发送图片消息
func (t *TelegramNotifier) sendPhotoMessage(ctx context.Context, chatID string, message *NotificationMessage) error {
	caption := fmt.Sprintf("*%s*\n\n%s", message.Title, message.Content)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", t.config.BotToken)

	requestBody := map[string]interface{}{
		"chat_id":    chatID,
		"photo":      message.Image,
		"caption":    caption,
		"parse_mode": "Markdown",
	}

	// 如果有URL，添加inline keyboard按钮
	if message.URL != "" {
		requestBody["reply_markup"] = map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{
						"text": "🔗 查看详情",
						"url":  message.URL,
					},
				},
			},
		}
	}

	return t.sendRequest(ctx, apiURL, requestBody)
}

// sendRequest 发送HTTP请求
func (t *TelegramNotifier) sendRequest(ctx context.Context, apiURL string, requestBody map[string]interface{}) error {
	var result TelegramResponse

	resp, err := t.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(&result).
		Post(apiURL)

	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	body := resp.Body()
	logger.Debug("telegram response: %s", string(body))

	if !resp.IsSuccess() {
		return fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode())
	}

	if !result.OK {
		return fmt.Errorf("发送消息失败: %s (错误代码: %d)", result.Description, result.ErrorCode)
	}

	return nil
}
