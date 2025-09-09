package server

import (
	"fmt"
	"net/http"

	"github.com/jianxcao/notify/internal/config"

	"github.com/gin-gonic/gin"
)

// findAppByID 根据 appID 查找应用配置和对应的 map key
func (s *HTTPServer) findAppByID(appID string) (config.NotificationApp, string, bool) {
	for mapKey, app := range s.config.NotificationApps {
		if app.AppID == appID {
			return app, mapKey, true
		}
	}
	return config.NotificationApp{}, "", false
}

// setupAdminRoutes 设置管理相关路由
func (s *HTTPServer) setupAdminRoutes(api *gin.RouterGroup) {
	// 管理接口（使用Basic认证）
	admin := api.Group("/admin", s.authMiddleware.BasicAuthMiddleware())
	{
		// 应用配置管理
		s.setupAppManagementRoutes(admin)

		// 模板管理
		s.setupTemplateManagementRoutes(admin)

		// 通知服务管理
		s.setupNotifierManagementRoutes(admin)

		// 插件管理
		s.setupPluginManagementRoutes(admin)
	}
}

// setupAppManagementRoutes 设置应用管理路由
func (s *HTTPServer) setupAppManagementRoutes(admin *gin.RouterGroup) {
	apps := admin.Group("/apps")
	{
		apps.GET("", s.handleGetApps)                // 获取所有应用
		apps.POST("", s.handleCreateApp)             // 创建新应用
		apps.GET("/:appid", s.handleGetAppConfig)    // 获取单个应用配置
		apps.PUT("/:appid", s.handleUpdateAppConfig) // 更新应用配置
		apps.DELETE("/:appid", s.handleDeleteApp)    // 删除应用
	}
}

// setupNotifierManagementRoutes 设置通知服务管理路由
func (s *HTTPServer) setupNotifierManagementRoutes(admin *gin.RouterGroup) {
	notifiers := admin.Group("/notifiers")
	{
		notifiers.GET("", s.handleGetNotifiers)                   // 获取所有通知服务
		notifiers.GET("/:notifier", s.handleGetNotifierConfig)    // 获取单个通知服务配置
		notifiers.PUT("/:notifier", s.handleUpdateNotifierConfig) // 更新通知服务配置
		notifiers.DELETE("/:notifier", s.handleDeleteNotifier)    // 删除通知服务
	}
}

// setupTemplateManagementRoutes 设置模板管理路由
func (s *HTTPServer) setupTemplateManagementRoutes(admin *gin.RouterGroup) {
	templates := admin.Group("/templates")
	{
		templates.GET("", s.handleGetTemplates)                  // 获取所有模板
		templates.GET("/:templateId", s.handleGetTemplate)       // 获取单个模板
		templates.POST("", s.handleCreateTemplate)               // 创建模板
		templates.PUT("/:templateId", s.handleUpdateTemplate)    // 更新模板
		templates.DELETE("/:templateId", s.handleDeleteTemplate) // 删除模板
	}
}

// handleGetApps 获取所有通知应用
func (s *HTTPServer) handleGetApps(c *gin.Context) {
	apps := s.app.GetNotificationApps()
	c.JSON(http.StatusOK, NewSuccessRes(apps))
}

// handleGetAppConfig 获取单个应用配置
func (s *HTTPServer) handleGetAppConfig(c *gin.Context) {
	appID := c.Param("appid")

	app, _, found := s.findAppByID(appID)
	if !found {
		c.JSON(http.StatusOK, NewErrorRes(APP_NOT_FOUND, fmt.Sprintf("应用 %s 不存在", appID)))
		return
	}

	c.JSON(http.StatusOK, NewSuccessRes(app))
}

// handleUpdateAppConfig 更新应用配置
func (s *HTTPServer) handleUpdateAppConfig(c *gin.Context) {
	appID := c.Param("appid")

	var updateReq config.NotificationApp
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
		return
	}

	// 确保 AppID 与路径参数一致
	updateReq.AppID = appID

	// 直接使用 updateReq 参数更新应用配置
	if err := s.configManager.UpdateAppConfig(updateReq); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(APP_CONFIG_ERROR, "更新应用配置失败"))
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()
	s.app.InitNotifiers()
	c.JSON(http.StatusOK, NewSuccessRes(updateReq))
}

// handleCreateApp 创建新应用
func (s *HTTPServer) handleCreateApp(c *gin.Context) {
	var createReq config.NotificationApp

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
		return
	}

	// 确保 AppID 字段被设置
	if createReq.AppID == "" {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "app_id 字段不能为空"))
		return
	}

	// 检查 AppID 是否已存在
	if _, _, found := s.findAppByID(createReq.AppID); found {
		c.JSON(http.StatusOK, NewErrorRes(APP_ALREADY_EXISTS, fmt.Sprintf("应用 %s 已存在", createReq.AppID)))
		return
	}

	// 使用ConfigManager创建应用，直接使用 AppID 作为 map key
	if err := s.configManager.CreateApp(createReq.AppID, createReq); err != nil {
		if err.Error() == fmt.Sprintf("应用 %s 已存在", createReq.AppID) {
			c.JSON(http.StatusOK, NewErrorRes(APP_ALREADY_EXISTS, err.Error()))
		} else {
			c.JSON(http.StatusOK, NewErrorRes(APP_CONFIG_ERROR, "创建应用失败"))
		}
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()
	s.app.InitNotifiers()
	c.JSON(http.StatusCreated, NewSuccessRes(createReq))
}

// handleDeleteApp 删除应用
func (s *HTTPServer) handleDeleteApp(c *gin.Context) {
	appID := c.Param("appid")

	// 检查应用是否存在并获取 mapKey
	_, mapKey, found := s.findAppByID(appID)
	if !found {
		c.JSON(http.StatusOK, NewErrorRes(APP_NOT_FOUND, fmt.Sprintf("应用 %s 不存在", appID)))
		return
	}

	// 使用ConfigManager删除应用（使用 mapKey）
	if err := s.configManager.DeleteApp(mapKey); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(APP_CONFIG_ERROR, "删除应用失败"))
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()
	s.app.InitNotifiers()
	c.JSON(http.StatusOK, NewSuccessRes(fmt.Sprintf("应用 %s 删除成功", appID)))
}

// ===== 模板管理接口处理函数 =====

// handleGetTemplates 获取所有模板
func (s *HTTPServer) handleGetTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, NewSuccessRes(s.config.Templates))
}

// handleGetTemplate 获取单个模板
func (s *HTTPServer) handleGetTemplate(c *gin.Context) {
	templateId := c.Param("templateId")

	template, exists := s.config.Templates[templateId]
	if !exists {
		c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_NOT_FOUND, fmt.Sprintf("模板 %s 不存在", templateId)))
		return
	}

	c.JSON(http.StatusOK, NewSuccessRes(template))
}

// handleCreateTemplate 创建模板
func (s *HTTPServer) handleCreateTemplate(c *gin.Context) {
	var createReq config.MessageTemplate

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
		return
	}

	// 确保 TemplateID 字段被设置
	if createReq.ID == "" {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "template_id 字段不能为空"))
		return
	}

	// 检查 TemplateID 是否已存在
	if _, exists := s.config.Templates[createReq.ID]; exists {
		c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_ALREADY_EXISTS, fmt.Sprintf("模板 %s 已存在", createReq.ID)))
		return
	}

	// 使用ConfigManager创建模板
	if err := s.configManager.CreateTemplate(createReq.ID, createReq); err != nil {
		if err.Error() == fmt.Sprintf("模板 %s 已存在", createReq.ID) {
			c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_ALREADY_EXISTS, err.Error()))
		} else {
			c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_CONFIG_ERROR, "创建模板失败"))
		}
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()

	c.JSON(http.StatusCreated, NewSuccessRes(createReq))
}

// handleUpdateTemplate 更新模板
func (s *HTTPServer) handleUpdateTemplate(c *gin.Context) {
	templateId := c.Param("templateId")

	// 检查模板是否存在
	if _, exists := s.config.Templates[templateId]; !exists {
		c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_NOT_FOUND, fmt.Sprintf("模板 %s 不存在", templateId)))
		return
	}

	var updateReq config.MessageTemplate
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(PARAM_ERROR, "解析请求失败"))
		return
	}

	// 更新配置
	newTemplates := make(map[string]config.MessageTemplate)
	for k, v := range s.config.Templates {
		newTemplates[k] = v
	}
	newTemplates[templateId] = updateReq

	if err := s.configManager.UpdateTemplatesConfig(newTemplates); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_CONFIG_ERROR, "更新模板配置失败"))
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()

	c.JSON(http.StatusOK, NewSuccessRes(updateReq))
}

// handleDeleteTemplate 删除模板
func (s *HTTPServer) handleDeleteTemplate(c *gin.Context) {
	templateId := c.Param("templateId")

	// 检查模板是否存在
	if _, exists := s.config.Templates[templateId]; !exists {
		c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_NOT_FOUND, fmt.Sprintf("模板 %s 不存在", templateId)))
		return
	}
	// 检查是否有应用在使用这个模板
	appsUsingTemplate := s.configManager.GetAppsUsingTemplate(templateId)
	if len(appsUsingTemplate) > 0 {
		c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_ALREADY_EXISTS, fmt.Sprintf("模板 %s 正在被以下应用使用，不能删除: %v", templateId, appsUsingTemplate)))
		return
	}

	if err := s.configManager.DeleteTemplate(templateId); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(TEMPLATE_CONFIG_ERROR, "删除模板失败"))
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()

	c.JSON(http.StatusOK, NewSuccessRes(fmt.Sprintf("模板 %s 删除成功", templateId)))
}

// ===== 通知服务管理接口处理函数 =====

// handleGetNotifiers 获取所有通知服务实例
func (s *HTTPServer) handleGetNotifiers(c *gin.Context) {
	c.JSON(http.StatusOK, NewSuccessRes(s.config.Notifiers))
}

// NotifierConfigResponse 通知服务配置响应结构体（驼峰命名）
type NotifierConfigResponse struct {
	InstanceName string                 `json:"instanceName"`
	Type         config.NotifiersType   `json:"type"`
	Enabled      bool                   `json:"enabled"`
	Config       map[string]interface{} `json:"config"`
}

// handleGetNotifierConfig 获取特定通知服务实例配置
func (s *HTTPServer) handleGetNotifierConfig(c *gin.Context) {
	instanceName := c.Param("notifier")

	notifierInstance, exists := s.config.Notifiers[instanceName]
	if !exists {
		c.JSON(http.StatusOK, NewErrorRes(NOTIFIER_NOT_FOUND, fmt.Sprintf("通知服务实例 %s 不存在", instanceName)))
		return
	}

	// 隐藏敏感信息
	safeConfig := make(map[string]interface{})
	for key, value := range notifierInstance.Config {
		switch key {
		case "secret", "bot_token", "access_token":
			safeConfig[key] = "***"
		default:
			safeConfig[key] = value
		}
	}

	responseData := NotifierConfigResponse{
		InstanceName: instanceName,
		Type:         notifierInstance.Type,
		Enabled:      notifierInstance.Enabled,
		Config:       safeConfig,
	}

	c.JSON(http.StatusOK, NewSuccessRes(responseData))
}

// handleUpdateNotifierConfig 更新通知服务实例配置
func (s *HTTPServer) handleUpdateNotifierConfig(c *gin.Context) {
	instanceName := c.Param("notifier")

	// 检查实例是否存在就是新增
	// if _, exists := s.config.Notifiers[instanceName]; !exists {
	// 	c.JSON(http.StatusNotFound, NewErrorRes(NOTIFIER_NOT_FOUND, fmt.Sprintf("通知服务实例 %s 不存在", instanceName)))
	// 	return
	// }

	var updateReq config.NotifierInstance
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorRes(PARAM_ERROR, "解析请求失败"))
		return
	}

	// 更新配置
	newNotifiers := make(map[string]config.NotifierInstance)
	for k, v := range s.config.Notifiers {
		newNotifiers[k] = v
	}
	newNotifiers[instanceName] = updateReq

	if err := s.configManager.UpdateNotifiersConfig(newNotifiers); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(NOTIFIER_CONFIG_ERROR, "更新通知服务实例配置失败"))
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()
	s.app.InitNotifiers()
	c.JSON(http.StatusOK, NewSuccessRes(updateReq))
}

// NotifierDeleteResponse 通知服务删除时如果有应用在使用的响应结构体
type NotifierDeleteResponse struct {
	Message   string   `json:"message"`
	UsingApps []string `json:"usingApps"`
}

// handleDeleteNotifier 删除通知服务实例
func (s *HTTPServer) handleDeleteNotifier(c *gin.Context) {
	instanceName := c.Param("notifier")

	// 检查实例是否存在
	if _, exists := s.config.Notifiers[instanceName]; !exists {
		c.JSON(http.StatusOK, NewErrorRes(NOTIFIER_NOT_FOUND, fmt.Sprintf("通知服务实例 %s 不存在", instanceName)))
		return
	}

	// 检查是否有应用在使用这个通知服务实例
	usingApps := s.configManager.GetAppsUsingNotifier(instanceName)
	if len(usingApps) > 0 {
		responseData := NotifierDeleteResponse{
			Message:   fmt.Sprintf("无法删除通知服务实例 %s，以下应用正在使用", instanceName),
			UsingApps: usingApps,
		}
		c.JSON(http.StatusOK, NewErrorRes(NOTIFIER_IN_USE, responseData.Message))
		return
	}

	if err := s.configManager.DeleteNotifier(instanceName); err != nil {
		c.JSON(http.StatusOK, NewErrorRes(NOTIFIER_CONFIG_ERROR, "删除通知服务实例失败"))
		return
	}

	// 更新内存中的配置
	s.config = s.configManager.GetConfig()
	s.app.InitNotifiers()
	c.JSON(http.StatusOK, NewSuccessRes(fmt.Sprintf("通知服务实例 %s 删除成功", instanceName)))
}
