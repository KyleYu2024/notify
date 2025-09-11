package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type NotifiersType string

var (
	WechatWorkAPPBot     NotifiersType = "wechatWorkAPPBot"
	WechatWorkWebhookBot NotifiersType = "wechatWorkWebhookBot"
	TelegramAppBot       NotifiersType = "telegramAppBot"
	DingTalkAppBot       NotifiersType = "dingTalkAppBot"
	FeishuAppBot         NotifiersType = "feishuAppBot"
)

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `yaml:"level" json:"level"`   // debug, info, warn, error
	Format string `yaml:"format" json:"format"` // json, text
}

// Config 主配置结构
type Config struct {
	Logger           *LoggerConfig               `yaml:"logger,omitempty" json:"logger,omitempty"`
	Notifiers        map[string]NotifierInstance `yaml:"notifiers" json:"notifiers"`
	Templates        map[string]MessageTemplate  `yaml:"templates" json:"templates"` // 消息模板配置
	NotificationApps map[string]NotificationApp  `yaml:"notification_apps" json:"notificationApps"`
}

// NotifierInstance 通知服务实例配置
type NotifierInstance struct {
	Type    NotifiersType          `yaml:"type" json:"type" binding:"required"`
	Enabled bool                   `yaml:"enabled" json:"enabled"`
	Config  map[string]interface{} `yaml:",inline" json:"config"`
}

// WechatWorkConfig 企业微信应用配置
type WechatWorkConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	CorpID  string `yaml:"corp_id" json:"corpId"`
	AgentID string `yaml:"agent_id" json:"agentId"`
	Secret  string `yaml:"secret" json:"secret"`
	Targets string `yaml:"targets" json:"targets"`
	Proxy   string `yaml:"proxy" json:"proxy"` // 代理服务器地址，格式: http://proxy.example.com:8080
}

// WechatWorkWebhookConfig 企业微信群机器人配置
type WechatWorkWebhookConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Key     string `yaml:"key" json:"key"`     // 群机器人的 key
	Proxy   string `yaml:"proxy" json:"proxy"` // 代理服务器地址，格式: http://proxy.example.com:8080
}

// TelegramConfig Telegram配置
type TelegramConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	BotToken string `yaml:"bot_token" json:"botToken"`
	ChatID   string `yaml:"chat_id" json:"chatId"` // 新增chatId字段
	Proxy    string `yaml:"proxy" json:"proxy"`    // 代理服务器地址，格式: http://proxy.example.com:8080
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	Enabled     bool   `yaml:"enabled" json:"enabled"`
	AccessToken string `yaml:"access_token" json:"accessToken"`
	Secret      string `yaml:"secret" json:"secret"`
	Targets     string `yaml:"targets" json:"targets"`
	Proxy       string `yaml:"proxy" json:"proxy"` // 代理服务器地址，格式: http://proxy.example.com:8080
}

// FeishuConfig 飞书配置
type FeishuConfig struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	AppID     string `yaml:"app_id" json:"appId"`         // 飞书应用ID
	AppSecret string `yaml:"app_secret" json:"appSecret"` // 飞书应用密钥

	Targets string `yaml:"targets" json:"targets"` // 默认发送目标(用户ID或群ID)
	Proxy   string `yaml:"proxy" json:"proxy"`     // 代理服务器地址，格式: http://proxy.example.com:8080
}

// NotificationApp 通知应用配置
type NotificationApp struct {
	AppID        string   `yaml:"app_id" json:"appId" binding:"required"`
	Name         string   `yaml:"name" json:"name"`
	Description  string   `yaml:"description" json:"description"`
	Enabled      bool     `yaml:"enabled" json:"enabled"`
	Notifiers    []string `yaml:"notifiers" json:"notifiers"`
	TemplateID   string   `yaml:"template_id" json:"templateId"`                 // 关联的模板ID
	PluginID     string   `yaml:"plugin_id,omitempty" json:"pluginId,omitempty"` // 关联的插件ID
	DefaultImage string   `yaml:"default_image" json:"defaultImage"`             // 默认图片URL
	Auth         *AppAuth `yaml:"auth,omitempty" json:"auth,omitempty"`          // 可选字段
}

// AppAuth 通知应用的认证配置
type AppAuth struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Token   string `yaml:"token" json:"token"`
}

// MessageTemplate 消息模板配置
type MessageTemplate struct {
	ID      string `yaml:"id" json:"id"`           // 模板ID
	Name    string `yaml:"name" json:"name"`       // 模板名称
	Title   string `yaml:"title" json:"title"`     // 标题
	Content string `yaml:"content" json:"content"` // 内容
	Image   string `yaml:"image" json:"image"`     // 图片
	URL     string `yaml:"url" json:"url"`         // 链接
	Targets string `yaml:"targets" json:"targets"` // 目标
}

// ConfigManager 配置管理器
type ConfigManager struct {
	configFile string
	config     *Config
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configFile string) *ConfigManager {
	return &ConfigManager{
		configFile: configFile,
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// Load 加载配置
func (cm *ConfigManager) Load() (*Config, error) {
	config, err := LoadConfig(cm.configFile)
	if err != nil {
		return nil, err
	}
	cm.config = config
	if cm.config == nil {
		cm.config = &Config{}
	}
	if cm.config.Notifiers == nil {
		cm.config.Notifiers = make(map[string]NotifierInstance)
	}
	if cm.config.Templates == nil {
		cm.config.Templates = make(map[string]MessageTemplate)
	}
	if cm.config.NotificationApps == nil {
		cm.config.NotificationApps = make(map[string]NotificationApp)
	}
	cm.Save()
	return cm.config, nil
}

// Save 保存配置
func (cm *ConfigManager) Save() error {
	if cm.config == nil {
		return fmt.Errorf("配置未初始化")
	}
	return SaveConfig(cm.config, cm.configFile)
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// DeleteNotifier 删除通知服务实例
func (cm *ConfigManager) DeleteNotifier(instanceName string) error {
	delete(cm.config.Notifiers, instanceName)
	return SaveConfig(cm.config, cm.configFile)
}

// CreateTemplate 创建新模板
func (cm *ConfigManager) CreateTemplate(templateID string, template MessageTemplate) error {
	if _, exists := cm.config.Templates[templateID]; exists {
		return fmt.Errorf("模板 %s 已存在", templateID)
	}

	cm.config.Templates[templateID] = template
	return SaveConfig(cm.config, cm.configFile)
}

// UpdateTemplatesConfig 更新模板配置
func (cm *ConfigManager) UpdateTemplatesConfig(templates map[string]MessageTemplate) error {
	cm.config.Templates = templates
	return SaveConfig(cm.config, cm.configFile)
}

// DeleteTemplate 删除模板
func (cm *ConfigManager) DeleteTemplate(templateID string) error {
	delete(cm.config.Templates, templateID)
	return SaveConfig(cm.config, cm.configFile)
}

// GetAppsUsingTemplate 获取使用指定模板的应用列表
func (cm *ConfigManager) GetAppsUsingTemplate(templateID string) []string {
	var apps []string
	for appID, app := range cm.config.NotificationApps {
		if app.TemplateID == templateID {
			apps = append(apps, appID)
		}
	}
	return apps
}

// UpdateApp 更新应用配置
func (cm *ConfigManager) UpdateApp(appName string, updates map[string]interface{}) error {
	if cm.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	app, exists := cm.config.NotificationApps[appName]
	if !exists {
		return fmt.Errorf("应用 %s 不存在", appName)
	}

	if name, ok := updates["name"].(string); ok {
		app.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		app.Description = description
	}
	if enabled, ok := updates["enabled"].(bool); ok {
		app.Enabled = enabled
	}
	if notifiers, ok := updates["notifiers"].([]interface{}); ok {
		stringNotifiers := make([]string, len(notifiers))
		for i, n := range notifiers {
			if str, ok := n.(string); ok {
				stringNotifiers[i] = str
			}
		}
		app.Notifiers = stringNotifiers
	}
	if templateID, ok := updates["template_id"].(string); ok {
		app.TemplateID = templateID
	}
	if pluginID, ok := updates["plugin_id"].(string); ok {
		app.PluginID = pluginID
	}
	if authData, ok := updates["auth"].(map[string]interface{}); ok {
		// 如果Auth字段为nil，先初始化
		if app.Auth == nil {
			app.Auth = &AppAuth{}
		}

		if enabled, ok := authData["enabled"].(bool); ok {
			app.Auth.Enabled = enabled
		}
		if token, ok := authData["token"].(string); ok {
			app.Auth.Token = token
		}
	}

	cm.config.NotificationApps[appName] = app
	return cm.Save()
}

// UpdateAppConfig 更新应用配置（直接使用 appConfig.AppID）
func (cm *ConfigManager) UpdateAppConfig(appConfig NotificationApp) error {
	if cm.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	if appConfig.AppID == "" {
		return fmt.Errorf("AppID 不能为空")
	}

	// 查找对应的 map key（在我们的设计中，map key 应该等于 AppID）
	var mapKey string
	found := false
	for key, app := range cm.config.NotificationApps {
		if app.AppID == appConfig.AppID {
			mapKey = key
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("应用 %s 不存在", appConfig.AppID)
	}

	cm.config.NotificationApps[mapKey] = appConfig
	return cm.Save()
}

// UpdateNotifiersConfig 更新通知服务配置
func (cm *ConfigManager) UpdateNotifiersConfig(notifiersConfig map[string]NotifierInstance) error {
	if cm.config == nil {
		return fmt.Errorf("配置未初始化")
	}
	cm.config.Notifiers = notifiersConfig
	return cm.Save()
}

// CreateApp 创建新应用
func (cm *ConfigManager) CreateApp(appName string, appConfig NotificationApp) error {
	if cm.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	if _, exists := cm.config.NotificationApps[appName]; exists {
		return fmt.Errorf("应用 %s 已存在", appName)
	}

	cm.config.NotificationApps[appName] = appConfig
	return cm.Save()
}

// DeleteApp 删除应用
func (cm *ConfigManager) DeleteApp(appName string) error {
	if cm.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	if _, exists := cm.config.NotificationApps[appName]; !exists {
		return fmt.Errorf("应用 %s 不存在", appName)
	}

	delete(cm.config.NotificationApps, appName)
	return cm.Save()
}

// GetAppsUsingNotifier 获取使用指定通知服务的应用列表
func (cm *ConfigManager) GetAppsUsingNotifier(notifierName string) []string {
	if cm.config == nil {
		return nil
	}

	appsUsingNotifier := []string{}
	for appName, appConfig := range cm.config.NotificationApps {
		for _, notifier := range appConfig.Notifiers {
			if notifier == notifierName {
				appsUsingNotifier = append(appsUsingNotifier, appName)
				break
			}
		}
	}
	return appsUsingNotifier
}

// GetAppsUsingPlugin 获取使用指定插件的应用列表
func (cm *ConfigManager) GetAppsUsingPlugin(pluginID string) []string {
	if cm.config == nil {
		return nil
	}

	var apps []string
	for appName, app := range cm.config.NotificationApps {
		if app.PluginID == pluginID {
			apps = append(apps, appName)
		}
	}
	return apps
}
