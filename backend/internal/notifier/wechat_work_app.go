package notifier

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jianxcao/notify/internal/config"

	"github.com/go-resty/resty/v2"
)

// WechatWorkNotifier 企业微信通知服务
type WechatWorkNotifier struct {
	config      config.WechatWorkConfig
	accessToken string
	client      *resty.Client
	baseURL     string
}

// NewWechatWorkNotifier 创建企业微信通知服务实例
func NewWechatWorkNotifier(cfg config.WechatWorkConfig) *WechatWorkNotifier {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(2 * time.Second)

	baseURL := "https://qyapi.weixin.qq.com"
	if cfg.Proxy != "" {
		baseURL = strings.TrimSuffix(cfg.Proxy, "/")
	}

	return &WechatWorkNotifier{
		config:  cfg,
		client:  client,
		baseURL: baseURL,
	}
}

// Name 返回服务名称
func (w *WechatWorkNotifier) Name() string {
	return string(config.WechatWorkAPPBot)
}

// IsEnabled 检查服务是否启用
func (w *WechatWorkNotifier) IsEnabled() bool {
	return w.config.Enabled
}

// Validate 验证配置
func (w *WechatWorkNotifier) Validate() error {
	if !w.config.Enabled {
		return nil
	}

	if w.config.CorpID == "" {
		return fmt.Errorf("企业微信 CorpID 不能为空")
	}
	if w.config.AgentID == "" {
		return fmt.Errorf("企业微信 AgentID 不能为空")
	}
	if w.config.Secret == "" {
		return fmt.Errorf("企业微信 Secret 不能为空")
	}

	return nil
}

// TokenResponse 访问令牌响应结构
type TokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// MessageResponse 消息发送响应结构
type MessageResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// getAccessToken 获取访问令牌
func (w *WechatWorkNotifier) getAccessToken(ctx context.Context) error {
	var result TokenResponse

	resp, err := w.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"corpid":     w.config.CorpID,
			"corpsecret": w.config.Secret,
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/cgi-bin/gettoken", w.baseURL))
	fmt.Println("body", resp.RawBody())
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode())
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("获取访问令牌失败: %s", result.ErrMsg)
	}

	w.accessToken = result.AccessToken
	return nil
}

// Send 发送通知消息
func (w *WechatWorkNotifier) Send(ctx context.Context, message *NotificationMessage, targets []string) error {
	if !w.config.Enabled {
		return fmt.Errorf("企业微信通知服务未启用")
	}
	if len(targets) == 0 && w.config.Targets != "" {
		targets = strings.Split(w.config.Targets, ",")
	}

	// 获取访问令牌
	if err := w.getAccessToken(ctx); err != nil {
		return fmt.Errorf("获取访问令牌失败: %w", err)
	}
	// TODO: 消息超过2048字符的处理
	// 如果有图片，发送图文消息，否则发送文本消息
	if message.Image != "" {
		return w.sendNewsMessage(ctx, message, targets)
	}

	return w.sendTextMessage(ctx, message, targets)
}

// sendTextMessage 发送文本消息
func (w *WechatWorkNotifier) sendTextMessage(ctx context.Context, message *NotificationMessage, targets []string) error {
	// 构建消息内容
	content := fmt.Sprintf("%s\n%s", message.Title, message.Content)

	// 构建发送消息的请求体
	requestBody := map[string]interface{}{
		"touser":  "@all", // 默认发送给所有人，可以根据targets参数调整
		"msgtype": "text",
		"agentid": w.config.AgentID,
		"text": map[string]string{
			"content": content,
		},
	}

	// 如果指定了目标用户
	if len(targets) > 0 {
		requestBody["touser"] = joinStrings(targets, "|")
	}

	return w.sendMessage(ctx, requestBody)
}

// sendNewsMessage 发送图文消息
func (w *WechatWorkNotifier) sendNewsMessage(ctx context.Context, message *NotificationMessage, targets []string) error {
	// 构建图文消息
	articles := []map[string]interface{}{
		{
			"title":       message.Title,
			"description": message.Content,
			"url":         message.URL, // 使用通知消息中的URL
			"picurl":      message.Image,
		},
	}

	requestBody := map[string]interface{}{
		"touser":  "@all",
		"msgtype": "news",
		"agentid": w.config.AgentID,
		"news": map[string]interface{}{
			"articles": articles,
		},
	}

	// 如果指定了目标用户
	if len(targets) > 0 {
		requestBody["touser"] = joinStrings(targets, "|")
	}

	return w.sendMessage(ctx, requestBody)
}

// sendMessage 发送消息到企业微信
func (w *WechatWorkNotifier) sendMessage(ctx context.Context, requestBody map[string]interface{}) error {
	var result MessageResponse

	resp, err := w.client.R().
		SetContext(ctx).
		SetQueryParam("access_token", w.accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(&result).
		Post(fmt.Sprintf("%s/cgi-bin/message/send", w.baseURL))

	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode())
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("发送消息失败: %s", result.ErrMsg)
	}

	return nil
}

// joinStrings 连接字符串切片
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	var result string
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
