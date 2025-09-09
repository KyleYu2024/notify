package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jianxcao/notify/internal/config"

	"github.com/go-resty/resty/v2"
)

// WechatWorkWebhookNotifier 企业微信群机器人通知服务
type WechatWorkWebhookNotifier struct {
	config  config.WechatWorkWebhookConfig
	client  *resty.Client
	baseURL string
}

// NewWechatWorkWebhookNotifier 创建企业微信群机器人通知服务实例
func NewWechatWorkWebhookNotifier(cfg config.WechatWorkWebhookConfig) *WechatWorkWebhookNotifier {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(2 * time.Second)

	baseURL := "https://qyapi.weixin.qq.com"
	if cfg.Proxy != "" {
		baseURL = strings.TrimSuffix(cfg.Proxy, "/")
	}

	return &WechatWorkWebhookNotifier{
		config:  cfg,
		client:  client,
		baseURL: baseURL,
	}
}

// Name 返回服务名称
func (w *WechatWorkWebhookNotifier) Name() string {
	return string(config.WechatWorkWebhookBot)
}

// IsEnabled 检查服务是否启用
func (w *WechatWorkWebhookNotifier) IsEnabled() bool {
	return w.config.Enabled
}

// Validate 验证配置是否有效
func (w *WechatWorkWebhookNotifier) Validate() error {
	if !w.config.Enabled {
		return fmt.Errorf("企业微信群机器人未启用")
	}

	if w.config.Key == "" {
		return fmt.Errorf("企业微信群机器人 Key 不能为空")
	}

	return nil
}

// Send 发送通知消息
func (w *WechatWorkWebhookNotifier) Send(ctx context.Context, message *NotificationMessage, targets []string) error {
	if err := w.Validate(); err != nil {
		return err
	}

	// 发送消息
	if message.Image != "" {
		return w.SendNewsdownMessage(ctx, message)
	}
	return w.sendMarkdownMessage(ctx, message)
}

// buildMessage 构建消息内容
func (w *WechatWorkWebhookNotifier) SendTextMessage(ctx context.Context, message *NotificationMessage) error {
	// 构建完整的 webhook URL
	webhookURL := fmt.Sprintf("%s/cgi-bin/webhook/send?key=%s", w.baseURL, w.config.Key)

	// 企业微信群机器人消息格式
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": message.Title + "\n" + message.Content,
		},
	}
	resp, err := w.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(webhookURL)

	if err != nil {
		return fmt.Errorf("发送企业微信群机器人消息失败: %w", err)
	}

	return w.checkResp(resp)
}

func (w *WechatWorkWebhookNotifier) SendNewsdownMessage(ctx context.Context, message *NotificationMessage) error {
	// 构建完整的 webhook URL
	webhookURL := fmt.Sprintf("%s/cgi-bin/webhook/send?key=%s", w.baseURL, w.config.Key)
	url := message.URL
	if url == "" {
		url = message.Image
	}
	// 企业微信群机器人消息格式
	payload := map[string]interface{}{
		"msgtype": "news",
		"news": map[string]interface{}{
			"articles": []map[string]interface{}{
				{
					"title":       message.Title,
					"description": message.Content,
					"url":         url,
					"picurl":      message.Image,
				},
			},
		},
	}

	resp, err := w.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(webhookURL)

	if err != nil {
		return fmt.Errorf("发送企业微信群机器人消息失败: %w", err)
	}

	return w.checkResp(resp)
}

// sendWebhookMessage 发送 webhook 消息
func (w *WechatWorkWebhookNotifier) sendMarkdownMessage(ctx context.Context, message *NotificationMessage) error {

	var content strings.Builder
	if message.Image != "" {
		content.WriteString(fmt.Sprintf("![%s](%s)\n\n", message.Title, message.Image))
	}
	// 添加标题
	if message.Title != "" {
		content.WriteString(fmt.Sprintf("**%s**\n\n", message.Title))
	}

	// 添加内容
	if message.Content != "" {
		content.WriteString(message.Content)
		content.WriteString("\n\n")
	}

	// 添加时间戳
	if message.Timestamp != "" {
		content.WriteString(fmt.Sprintf("⏰ %s", message.Timestamp))
	}

	// 构建完整的 webhook URL
	webhookURL := fmt.Sprintf("%s/cgi-bin/webhook/send?key=%s", w.baseURL, w.config.Key)

	// 企业微信群机器人消息格式
	payload := map[string]interface{}{
		"msgtype": "markdown_v2",
		"markdown_v2": map[string]interface{}{
			"content": content.String(),
		},
	}

	resp, err := w.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(webhookURL)

	if err != nil {
		return fmt.Errorf("发送企业微信群机器人消息失败: %w", err)
	}

	return w.checkResp(resp)
}

func (w *WechatWorkWebhookNotifier) checkResp(resp *resty.Response) error {
	if resp.StatusCode() != 200 {
		return fmt.Errorf("企业微信群机器人API返回错误状态码: %d, 响应: %s", resp.StatusCode(), resp.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("解析企业微信群机器人响应失败: %w", err)
	}

	// 检查是否成功
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		errmsg, _ := result["errmsg"].(string)
		return fmt.Errorf("企业微信群机器人返回错误: errcode=%v, errmsg=%s", errcode, errmsg)
	}
	return nil
}
