package notifier

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jianxcao/notify/backend/pkg/config"
	"github.com/jianxcao/notify/backend/pkg/utils"

	"github.com/go-resty/resty/v2"
)

// DingTalkNotifier 钉钉通知服务
type DingTalkNotifier struct {
	config config.DingTalkConfig
	client *resty.Client
}

// DingTalkResponse 钉钉API响应结构
type DingTalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// NewDingTalkNotifier 创建钉钉通知服务实例
func NewDingTalkNotifier(cfg config.DingTalkConfig) *DingTalkNotifier {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(2 * time.Second)

	// 如果配置了代理，设置代理
	if cfg.Proxy != "" {
		client.SetProxy(cfg.Proxy)
	}

	return &DingTalkNotifier{
		config: cfg,
		client: client,
	}
}

// Name 返回服务名称
func (d *DingTalkNotifier) Name() string {
	return string(config.DingTalkAppBot)
}

// IsEnabled 检查服务是否启用
func (d *DingTalkNotifier) IsEnabled() bool {
	return d.config.Enabled
}

// Validate 验证配置
func (d *DingTalkNotifier) Validate() error {
	if !d.config.Enabled {
		return nil
	}

	if d.config.AccessToken == "" {
		return fmt.Errorf("钉钉 Access Token 不能为空")
	}

	return nil
}

// generateSign 生成签名
func (d *DingTalkNotifier) generateSign(timestamp int64) (string, error) {
	if d.config.Secret == "" {
		return "", nil
	}

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, d.config.Secret)
	h := hmac.New(sha256.New, []byte(d.config.Secret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return url.QueryEscape(signature), nil
}

// Send 发送通知消息
func (d *DingTalkNotifier) Send(ctx context.Context, message *NotificationMessage, targets []string) error {
	if !d.config.Enabled {
		return fmt.Errorf("钉钉通知服务未启用")
	}
	if len(targets) == 0 && d.config.Targets != "" {
		targets = strings.Split(d.config.Targets, ",")
	}

	// 构建查询参数
	queryParams := map[string]string{
		"access_token": d.config.AccessToken,
	}

	// 如果配置了密钥，添加签名
	if d.config.Secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		sign, err := d.generateSign(timestamp)
		if err != nil {
			return fmt.Errorf("生成签名失败: %w", err)
		}
		queryParams["timestamp"] = fmt.Sprintf("%d", timestamp)
		queryParams["sign"] = sign
	}

	// 根据是否有图片选择不同的消息类型
	var requestBody map[string]interface{}
	// if message.Image != "" {
	// 	requestBody = d.buildFeedCardMessage(message, targets)
	// } else {
	// 	requestBody = d.buildMarkdownMessage(message, targets)
	// }
	//全部用 markdown
	requestBody = d.buildMarkdownMessage(message, targets)

	return d.sendMessage(ctx, queryParams, requestBody)
}

// buildMarkdownMessage 构建Markdown消息
func (d *DingTalkNotifier) buildMarkdownMessage(message *NotificationMessage, targets []string) map[string]interface{} {
	content := message.Content
	if message.Image != "" {
		content = fmt.Sprintf("![](%s)\n\n%s", message.Image, message.Content)
	}

	requestBody := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": message.Title,
			"text":  content,
		},
	}
	// 使用正则表达式检测 target 是不是一个有效的手机号，是手机号的放在一起，其他类型的放在一起
	phoneTargets := []string{}
	otherTargets := []string{}
	for _, target := range targets {
		if utils.IsMobilePhone(target) {
			phoneTargets = append(phoneTargets, target)
		} else {
			otherTargets = append(otherTargets, target)
		}
	}
	// 如果指定了@用户
	if len(targets) > 0 {
		requestBody["at"] = map[string]interface{}{
			"atMobiles": phoneTargets,
			"atUserIds": otherTargets,
			"isAtAll":   false,
		}
	}

	return requestBody
}

// buildFeedCardMessage 构建FeedCard消息（支持图片）
func (d *DingTalkNotifier) buildFeedCardMessage(message *NotificationMessage, targets []string) map[string]interface{} {
	links := []map[string]interface{}{
		{
			"title":      message.Title,
			"messageURL": message.URL, // 使用通知消息中的URL
			"picURL":     message.Image,
		},
	}

	requestBody := map[string]interface{}{
		"msgtype": "feedCard",
		"feedCard": map[string]interface{}{
			"links": links,
		},
	}

	// FeedCard类型不支持@功能，如果需要@用户，可以先发送一个文本消息
	if len(targets) > 0 {
		// 这里可以考虑先发送一个@消息，再发送图片卡片
		// 为了简化，暂时忽略@功能
	}

	return requestBody
}

// sendMessage 发送消息到钉钉
func (d *DingTalkNotifier) sendMessage(ctx context.Context, queryParams map[string]string, requestBody map[string]interface{}) error {
	var result DingTalkResponse

	resp, err := d.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(&result).
		Post("https://oapi.dingtalk.com/robot/send")

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
