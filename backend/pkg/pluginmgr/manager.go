package pluginmgr

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime"
	"sort"
	"time"

	"github.com/jianxcao/notify/backend/pkg/logger"
	"github.com/jianxcao/notify/backend/pkg/pluginsdk"
)

// PluginConfig 插件配置结构
type PluginConfig struct {
	// 插件元信息
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author,omitempty"`

	// 插件设置
	Settings map[string]any `json:"settings,omitempty"`

	// UI配置
	UI *pluginsdk.UIConfig `json:"ui,omitempty"`

	// 插件状态
	Enabled    bool   `json:"enabled"`
	TestData   any    `json:"test_data,omitempty"`
	ConfigFile string `json:"-"` // plugin.json 文件路径
}

// LoadedPlugin 已加载的插件实例
type LoadedPlugin struct {
	Config   *PluginConfig
	Instance pluginsdk.Plugin
	LoadedAt time.Time
}

// Manager 插件管理器
type Manager struct {
	// 插件根目录
	pluginsDir string

	// 已加载的插件 map[pluginID]*LoadedPlugin
	plugins map[string]*LoadedPlugin
}

// NewManager 创建插件管理器
func NewManager(pluginsDir string) *Manager {
	return &Manager{
		pluginsDir: pluginsDir,
		plugins:    make(map[string]*LoadedPlugin),
	}
}

// findPluginFile 根据当前系统查找合适的插件文件
func (m *Manager) findPluginFile(pluginDir string) (string, error) {
	// 构建可能的插件文件名列表，按优先级排序
	possibleFiles := []string{
		fmt.Sprintf("plugin-%s-%s.so", runtime.GOOS, runtime.GOARCH), // 当前系统架构
		"plugin.so", // 通用备用文件
	}

	// 按优先级查找文件
	for _, filename := range possibleFiles {
		fullPath := filepath.Join(pluginDir, filename)
		if _, err := os.Stat(fullPath); err == nil {
			logger.Error(fmt.Sprintf("找到适用于系统 %s-%s 的插件文件", runtime.GOOS, runtime.GOARCH), "路径", fullPath)
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("未找到适用于系统 %s-%s 的插件文件", runtime.GOOS, runtime.GOARCH)
}

// LoadAll 加载所有插件
func (m *Manager) LoadAll() error {
	// 检查插件目录是否存在
	if _, err := os.Stat(m.pluginsDir); os.IsNotExist(err) {
		logger.Warn(fmt.Sprintf("插件目录不存在，跳过插件加载: %s", m.pluginsDir))
		return nil
	}

	// 读取插件目录
	entries, err := os.ReadDir(m.pluginsDir)
	if err != nil {
		return fmt.Errorf("读取插件目录失败: %w", err)
	}

	// 遍历每个子目录
	loadCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginDir := filepath.Join(m.pluginsDir, entry.Name())
		if err := m.loadOne(pluginDir); err != nil {
			logger.Error(fmt.Sprintf("加载插件失败 [%s]: %v", entry.Name(), err))
			continue
		}

		loadCount++
		logger.Info(fmt.Sprintf("成功加载插件: %s", entry.Name()))
	}

	logger.Info(fmt.Sprintf("插件加载完成，共加载 %d 个插件", loadCount))
	return nil
}

// loadOne 加载单个插件
func (m *Manager) loadOne(pluginDir string) error {
	// 读取插件配置文件
	configFile := filepath.Join(pluginDir, "setting.json")
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("读取插件配置文件失败: %w", err)
	}

	// 解析插件配置
	var config PluginConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("解析插件配置失败: %w", err)
	}

	// 验证必填字段
	if config.ID == "" {
		return fmt.Errorf("插件ID不能为空")
	}
	if config.Name == "" {
		return fmt.Errorf("插件名称不能为空")
	}
	// 设置完整路径
	config.ConfigFile = configFile

	// 查找合适的插件文件
	soFile, err := m.findPluginFile(pluginDir)
	if err != nil {
		return fmt.Errorf("查找插件文件失败: %w", err)
	}

	// 打开插件
	p, err := plugin.Open(soFile)
	if err != nil {
		return fmt.Errorf("打开插件文件失败: %w", err)
	}

	// 查找插件构造函数
	sym, err := p.Lookup("NewPlugin")
	if err != nil {
		return fmt.Errorf("插件构造函数 NewPlugin 不存在: %w", err)
	}

	// 使用反射调用函数，避免类型断言问题
	funcValue := reflect.ValueOf(sym)
	funcType := funcValue.Type()

	// 验证函数签名
	if funcType.NumIn() != 0 || funcType.NumOut() != 2 {
		return fmt.Errorf("插件构造函数签名错误")
	}

	// 调用函数
	results := funcValue.Call([]reflect.Value{})

	// 检查错误
	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return fmt.Errorf("创建插件实例失败: %w", err)
	}

	// 获取插件实例并进行类型断言
	instanceInterface := results[0].Interface()
	instance, err := pluginsdk.WrapPlugin(instanceInterface)
	if err != nil {
		return fmt.Errorf("插件实例类型错误，期望 pluginsdk.Plugin，实际 %T", instanceInterface)
	}

	// 验证插件ID一致性
	if instance.ID() != config.ID {
		return fmt.Errorf("插件ID不匹配: 配置=%s, 实例=%s", config.ID, instance.ID())
	}

	// 合并设置：默认设置 + 配置文件设置
	finalSettings := make(map[string]any)

	// 1. 插件默认设置
	if defaultSettings := instance.DefaultSettings(); defaultSettings != nil {
		for k, v := range defaultSettings {
			finalSettings[k] = v
		}
	}

	// 2. 配置文件中的设置（优先级更高）
	if config.Settings != nil {
		for k, v := range config.Settings {
			finalSettings[k] = v
		}
	}

	config.Settings = finalSettings

	// 创建已加载插件实例
	loadedPlugin := &LoadedPlugin{
		Config:   &config,
		Instance: instance,
		LoadedAt: time.Now(),
	}

	// 注册插件
	m.plugins[config.ID] = loadedPlugin

	return nil
}

// GetPlugin 获取插件实例
func (m *Manager) GetPlugin(pluginID string) (*LoadedPlugin, bool) {
	plugin, exists := m.plugins[pluginID]
	return plugin, exists
}

// GetAllPlugins 获取所有已加载的插件
func (m *Manager) GetAllPlugins() map[string]*LoadedPlugin {
	return m.plugins
}

// ProcessChain 处理插件链
func (m *Manager) ProcessChain(ctx context.Context, pluginID string, input map[string]any) (*pluginsdk.Output, error) {
	// 获取插件
	loadedPlugin, exists := m.GetPlugin(pluginID)
	if !exists {
		return nil, fmt.Errorf("插件不存在: %s", pluginID)
	}

	// 检查插件是否启用
	if !loadedPlugin.Config.Enabled {
		return nil, fmt.Errorf("插件未启用: %s", pluginID)
	}

	// 执行插件处理
	output, err := loadedPlugin.Instance.Process(ctx, input, loadedPlugin.Config.Settings)
	if err != nil {
		return nil, fmt.Errorf("插件处理失败 [%s]: %w", pluginID, err)
	}

	// 设置元数据
	if output.Meta == nil {
		output.Meta = &pluginsdk.MetaData{}
	}
	output.Meta.PluginID = pluginID
	output.Meta.Req = input
	output.Meta.ProcessedAt = time.Now().Format(time.RFC3339)

	return output, nil
}

// UpdatePluginConfig 更新插件配置文件
func (m *Manager) UpdatePluginConfig(pluginID string, updates map[string]any) error {
	loadedPlugin, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("插件不存在: %s", pluginID)
	}

	// 读取当前配置
	config, err := m.loadPluginConfig(loadedPlugin.Config.ConfigFile)
	if err != nil {
		return fmt.Errorf("读取插件配置失败: %w", err)
	}

	// 更新配置
	if settings, ok := updates["settings"].(map[string]any); ok {
		config.Settings = settings
	}
	if enabled, ok := updates["enabled"].(bool); ok {
		config.Enabled = enabled
		loadedPlugin.Config.Enabled = enabled
	}

	// 保存配置文件
	if err := m.savePluginConfig(loadedPlugin.Config.ConfigFile, config); err != nil {
		return err
	}

	// 重新合并设置：默认设置 + 配置文件设置
	finalSettings := make(map[string]any)

	// 1. 插件默认设置
	if defaultSettings := loadedPlugin.Instance.DefaultSettings(); defaultSettings != nil {
		for k, v := range defaultSettings {
			finalSettings[k] = v
		}
	}

	// 2. 配置文件中的设置（优先级更高）
	if config.Settings != nil {
		for k, v := range config.Settings {
			finalSettings[k] = v
		}
	}

	// 更新内存中的设置
	loadedPlugin.Config.Settings = finalSettings

	return nil
}

// loadPluginConfig 加载插件配置文件
func (m *Manager) loadPluginConfig(configFile string) (*PluginConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config PluginConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// savePluginConfig 保存插件配置文件
func (m *Manager) savePluginConfig(configFile string, config *PluginConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	return os.WriteFile(configFile, data, 0644)
}

// IsPluginEnabled 检查插件是否启用
func (m *Manager) IsPluginEnabled(pluginID string) bool {
	if loadedPlugin, exists := m.plugins[pluginID]; exists {
		return loadedPlugin.Config.Enabled
	}
	return false
}

// GetPluginList 获取插件列表（用于API返回）
func (m *Manager) GetPluginList() []PluginInfo {
	var list []PluginInfo

	for _, loadedPlugin := range m.plugins {
		info := PluginInfo{
			ID:          loadedPlugin.Config.ID,
			Name:        loadedPlugin.Config.Name,
			Version:     loadedPlugin.Config.Version,
			Description: loadedPlugin.Config.Description,
			Author:      loadedPlugin.Config.Author,
			Enabled:     loadedPlugin.Config.Enabled,
			LoadedAt:    loadedPlugin.LoadedAt.UnixMilli(),
			UI:          loadedPlugin.Config.UI,
			Settings:    loadedPlugin.Config.Settings,
			TestData:    loadedPlugin.Config.TestData,
		}

		list = append(list, info)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}

// PluginInfo 插件信息结构（用于API返回）
type PluginInfo struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	Description string              `json:"description"`
	Author      string              `json:"author,omitempty"`
	Enabled     bool                `json:"enabled"`
	LoadedAt    int64               `json:"loadedAt"`
	UI          *pluginsdk.UIConfig `json:"ui"`
	Settings    map[string]any      `json:"settings"`
	TestData    any                 `json:"test_data"`
}
