package server

import (
	"fmt"
	"net/http"

	"github.com/jianxcao/notify/internal/pluginmgr"

	"github.com/gin-gonic/gin"
)

// setupPluginManagementRoutes 设置插件管理路由
func (s *HTTPServer) setupPluginManagementRoutes(admin *gin.RouterGroup) {
	plugins := admin.Group("/plugins")
	{
		plugins.GET("", s.handleGetPlugins)                          // 获取所有插件
		plugins.GET("/:pluginId", s.handleGetPlugin)                 // 获取单个插件信息
		plugins.PUT("/:pluginId/config", s.handleUpdatePluginConfig) // 更新插件配置
		plugins.PUT("/:pluginId/enable", s.handleEnablePlugin)       // 启用插件
		plugins.PUT("/:pluginId/disable", s.handleDisablePlugin)     // 禁用插件
		plugins.POST("/:pluginId/test", s.handleTestPlugin)          // 测试插件
	}
}

// handleGetPlugins 获取所有插件
func (s *HTTPServer) handleGetPlugins(c *gin.Context) {
	pluginManager := s.app.GetPluginManager()
	if pluginManager == nil {
		c.JSON(http.StatusOK, NewErrorRes(SERVER_ERROR, "插件管理器未初始化"))
		return
	}

	plugins := pluginManager.GetPluginList()
	c.JSON(http.StatusOK, NewSuccessRes(plugins))
}

// handleGetPlugin 获取单个插件信息
func (s *HTTPServer) handleGetPlugin(c *gin.Context) {
	pluginID := c.Param("pluginId")

	pluginManager := s.app.GetPluginManager()
	if pluginManager == nil {
		c.JSON(http.StatusOK, NewErrorRes(SERVER_ERROR, "插件管理器未初始化"))
		return
	}

	loadedPlugin, exists := pluginManager.GetPlugin(pluginID)
	if !exists {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_NOT_FOUND, fmt.Sprintf("插件 %s 不存在", pluginID)))
		return
	}

	// 构建插件信息
	pluginInfo := pluginmgr.PluginInfo{
		ID:          loadedPlugin.Config.ID,
		Name:        loadedPlugin.Config.Name,
		Version:     loadedPlugin.Config.Version,
		Description: loadedPlugin.Config.Description,
		Author:      loadedPlugin.Config.Author,
		Enabled:     loadedPlugin.Config.Enabled,
		LoadedAt:    loadedPlugin.LoadedAt.Format("2006-01-02T15:04:05Z07:00"),
		UI:          loadedPlugin.Config.UI,
		Settings:    loadedPlugin.Config.Settings,
	}

	c.JSON(http.StatusOK, NewSuccessRes(pluginInfo))
}

// handleUpdatePluginConfig 更新插件配置
func (s *HTTPServer) handleUpdatePluginConfig(c *gin.Context) {
	pluginID := c.Param("pluginId")

	var updateReq struct {
		Settings map[string]any `json:"settings"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
		return
	}

	pluginManager := s.app.GetPluginManager()
	if pluginManager == nil {
		c.JSON(http.StatusOK, NewErrorRes(SERVER_ERROR, "插件管理器未初始化"))
		return
	}

	// 检查插件是否存在
	_, exists := pluginManager.GetPlugin(pluginID)
	if !exists {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_NOT_FOUND, fmt.Sprintf("插件 %s 不存在", pluginID)))
		return
	}

	// 直接更新插件配置文件
	if err := pluginManager.UpdatePluginConfig(pluginID, map[string]any{
		"settings": updateReq.Settings,
	}); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_CONFIG_ERROR, "更新插件配置失败"))
		return
	}

	c.JSON(http.StatusOK, NewSuccessRes("插件配置更新成功"))
}

// handleEnablePlugin 启用插件
func (s *HTTPServer) handleEnablePlugin(c *gin.Context) {
	pluginID := c.Param("pluginId")

	pluginManager := s.app.GetPluginManager()
	if pluginManager == nil {
		c.JSON(http.StatusOK, NewErrorRes(SERVER_ERROR, "插件管理器未初始化"))
		return
	}

	// 检查插件是否存在
	_, exists := pluginManager.GetPlugin(pluginID)
	if !exists {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_NOT_FOUND, fmt.Sprintf("插件 %s 不存在", pluginID)))
		return
	}

	// 直接更新插件配置文件
	if err := pluginManager.UpdatePluginConfig(pluginID, map[string]any{
		"enabled": true,
	}); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_CONFIG_ERROR, "启用插件失败"))
		return
	}

	c.JSON(http.StatusOK, NewSuccessRes("插件启用成功"))
}

// handleDisablePlugin 禁用插件
func (s *HTTPServer) handleDisablePlugin(c *gin.Context) {
	pluginID := c.Param("pluginId")

	pluginManager := s.app.GetPluginManager()
	if pluginManager == nil {
		c.JSON(http.StatusOK, NewErrorRes(SERVER_ERROR, "插件管理器未初始化"))
		return
	}

	// 检查插件是否存在
	_, exists := pluginManager.GetPlugin(pluginID)
	if !exists {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_NOT_FOUND, fmt.Sprintf("插件 %s 不存在", pluginID)))
		return
	}

	// 检查是否有应用在使用该插件
	appsUsingPlugin := s.configManager.GetAppsUsingPlugin(pluginID)
	if len(appsUsingPlugin) > 0 {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_IN_USE, fmt.Sprintf("插件正在被以下应用使用，无法禁用: %v", appsUsingPlugin)))
		return
	}

	// 直接更新插件配置文件
	if err := pluginManager.UpdatePluginConfig(pluginID, map[string]any{
		"enabled": false,
	}); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_CONFIG_ERROR, "禁用插件失败"))
		return
	}

	c.JSON(http.StatusOK, NewSuccessRes("插件禁用成功"))
}

// handleTestPlugin 测试插件
func (s *HTTPServer) handleTestPlugin(c *gin.Context) {
	pluginID := c.Param("pluginId")

	var testReq struct {
		Input map[string]any `json:"input"`
	}

	if err := c.ShouldBindJSON(&testReq); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
		return
	}

	pluginManager := s.app.GetPluginManager()
	if pluginManager == nil {
		c.JSON(http.StatusOK, NewErrorRes(SERVER_ERROR, "插件管理器未初始化"))
		return
	}

	// 检查插件是否存在
	_, exists := pluginManager.GetPlugin(pluginID)
	if !exists {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_NOT_FOUND, fmt.Sprintf("插件 %s 不存在", pluginID)))
		return
	}

	// 执行插件处理
	output, err := pluginManager.ProcessChain(c.Request.Context(), pluginID, testReq.Input)
	if err != nil {
		c.JSON(http.StatusOK, NewErrorRes(PLUGIN_PROCESS_ERROR, fmt.Sprintf("插件处理失败: %v", err)))
		return
	}

	c.JSON(http.StatusOK, NewSuccessRes(map[string]any{
		"pluginId": pluginID,
		"input":    testReq.Input,
		"output":   output,
	}))
}
