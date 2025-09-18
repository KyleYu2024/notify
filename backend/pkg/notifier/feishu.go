package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jianxcao/notify/backend/pkg/config"
	"github.com/jianxcao/notify/backend/pkg/logger"

	"github.com/go-resty/resty/v2"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// FeishuNotifier 飞书通知服务
type FeishuNotifier struct {
	config     config.FeishuConfig
	larkClient *lark.Client // 飞书官方SDK客户端
}

// 不再需要额外的响应结构，直接使用官方SDK的响应

// NewFeishuNotifier 创建飞书通知服务实例
func NewFeishuNotifier(cfg config.FeishuConfig) *FeishuNotifier {
	notifier := &FeishuNotifier{
		config: cfg,
	}

	// 初始化飞书官方SDK客户端
	if cfg.AppID != "" && cfg.AppSecret != "" {
		notifier.larkClient = lark.NewClient(cfg.AppID, cfg.AppSecret)
	}

	return notifier
}

// Name 返回服务名称
func (f *FeishuNotifier) Name() string {
	return string(config.FeishuAppBot)
}

// IsEnabled 检查服务是否启用
func (f *FeishuNotifier) IsEnabled() bool {
	return f.config.Enabled
}

// Validate 验证配置
func (f *FeishuNotifier) Validate() error {
	if !f.config.Enabled {
		return nil
	}

	// 需要配置应用ID和密钥
	if f.config.AppID == "" || f.config.AppSecret == "" {
		return fmt.Errorf("飞书配置不完整：需要配置AppID和AppSecret")
	}

	return nil
}

// Send 发送通知消息
func (f *FeishuNotifier) Send(ctx context.Context, message *NotificationMessage, targets []string) error {
	if !f.config.Enabled {
		return fmt.Errorf("飞书通知服务未启用")
	}

	if len(targets) == 0 && f.config.Targets != "" {
		targets = strings.Split(f.config.Targets, ",")
	}

	// 使用飞书官方API发送消息
	if f.larkClient == nil {
		return fmt.Errorf("飞书客户端未初始化，请检查AppID和AppSecret配置")
	}

	return f.sendAPIMessage(ctx, message, targets)
}

// sendAPIMessage 通过API发送消息（应用机器人）
func (f *FeishuNotifier) sendAPIMessage(ctx context.Context, message *NotificationMessage, targets []string) error {
	// 如果没有指定目标，尝试发送到配置的默认目标
	if len(targets) == 0 {
		return fmt.Errorf("未指定消息发送目标")
	}

	// 为每个目标发送消息
	for _, target := range targets {
		target = strings.TrimSpace(target)
		if target == "" {
			continue
		}

		// 构建消息内容
		content := f.buildAPIMessageContent(message)

		// 判断目标类型并设置接收者ID类型
		receiveIdType := f.getReceiveIdType(target)

		// 创建请求
		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType(receiveIdType).
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(target).
				MsgType("post"). // 使用rich_text支持更丰富的格式
				Content(content).
				Build()).
			Build()

		// 发起请求
		resp, err := f.larkClient.Im.V1.Message.Create(ctx, req)
		if err != nil {
			return fmt.Errorf("发送消息到 %s 失败: %w", target, err)
		}

		// 检查响应
		if !resp.Success() {
			logger.Error("resp", resp.Err)
			logger.Error("err", resp.Code)
			logger.Error("err", resp.Msg)
			return fmt.Errorf("发送消息到 %s 失败: %s", target, resp.Msg)
		}
	}

	return nil
}

// buildAPIMessageContent 构建API消息内容（rich_text格式）
func (f *FeishuNotifier) buildAPIMessageContent(message *NotificationMessage) string {
	// 构建富文本内容
	content := map[string]interface{}{
		"zh_cn": map[string]interface{}{
			"title":   message.Title,
			"content": f.buildRichTextElements(message),
		},
	}

	// 转换为JSON字符串
	contentBytes, _ := json.Marshal(content)
	return string(contentBytes)
}

func (f *FeishuNotifier) buildImageElement(message *NotificationMessage) (string, error) {
	image, err := resty.New().R().SetDoNotParseResponse(true).Get(message.Image)
	if err != nil {
		return "", err
	}
	req := larkim.NewCreateImageReqBuilder().
		Body(larkim.NewCreateImageReqBodyBuilder().
			ImageType("message").
			Image(image.RawBody()).
			Build()).
		Build()

	resp, err := f.larkClient.Im.Image.Create(context.Background(), req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", resp.CodeError
	}

	return *resp.Data.ImageKey, nil
}

// buildRichTextElements 构建富文本元素
func (f *FeishuNotifier) buildRichTextElements(message *NotificationMessage) [][]map[string]interface{} {
	elements := [][]map[string]interface{}{}
	if message.Image != "" {
		imageKey, err := f.buildImageElement(message)
		if err != nil {
			logger.Error("buildImageElement", err)
		}
		elements = append(elements, []map[string]interface{}{
			{
				"tag":       "img",
				"image_key": imageKey,
			},
		})
	}
	// 添加内容行
	if message.Content != "" {
		// 按行分割内容
		content := []map[string]interface{}{
			{
				"tag":  "md",
				"text": message.Content,
			},
		}
		elements = append(elements, content)
	}

	// 添加时间戳
	if message.Timestamp != "" {
		timeElement := []map[string]interface{}{
			{
				"tag":  "text",
				"text": "时间: " + message.Timestamp,
			},
		}
		elements = append(elements, timeElement)
	}

	// 添加链接
	if message.URL != "" {
		linkElement := []map[string]interface{}{
			{
				"tag":  "a",
				"text": "查看详情",
				"href": message.URL,
			},
		}
		elements = append(elements, linkElement)
	}

	return elements
}

// getReceiveIdType 根据目标格式判断接收者ID类型
// 根据飞书官方文档：https://open.feishu.cn/document/server-docs/im-v1/message/create
func (f *FeishuNotifier) getReceiveIdType(target string) string {
	// 根据ID前缀判断类型
	if strings.HasPrefix(target, "ou_") {
		return "open_id" // Open ID: 用户在某个应用中的身份标识
	} else if strings.HasPrefix(target, "on_") {
		return "union_id" // Union ID: 用户在应用开发商下的统一身份标识
	} else if strings.HasPrefix(target, "oc_") {
		return "chat_id" // Chat ID: 群聊标识
	} else if strings.Contains(target, "@") {
		return "email" // Email: 用户邮箱
	} else {
		return "user_id" // User ID: 用户在租户内的身份标识（默认）
	}
}
