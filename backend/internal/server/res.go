package server

type BaseRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewBaseRes(code int, msg string, data interface{}) *BaseRes {
	return &BaseRes{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewSuccessRes(data interface{}) *BaseRes {
	return NewBaseRes(SUCCESS_CODE, "success", data)
}

func NewErrorRes(code int, msg string) *BaseRes {
	return NewBaseRes(code, msg, nil)
}

// 状态码定义
const (
	SUCCESS_CODE = 0 // 成功

	// 通用错误码 (1000-1999)
	PARAM_ERROR      = 1001 // 参数错误
	AUTH_ERROR       = 1002 // 认证失败
	PERMISSION_ERROR = 1003 // 权限不足
	NOT_FOUND_ERROR  = 1004 // 资源不存在
	CONFLICT_ERROR   = 1005 // 资源冲突
	VALIDATION_ERROR = 1006 // 数据验证失败

	// 应用相关错误码 (2000-2999)
	APP_NOT_FOUND      = 2001 // 应用不存在
	APP_DISABLED       = 2002 // 应用未启用
	APP_ALREADY_EXISTS = 2003 // 应用已存在
	APP_CONFIG_ERROR   = 2004 // 应用配置错误

	// 模板相关错误码 (3000-3999)
	TEMPLATE_NOT_FOUND      = 3001 // 模板不存在
	TEMPLATE_ALREADY_EXISTS = 3002 // 模板已存在
	TEMPLATE_CONFIG_ERROR   = 3003 // 模板配置错误

	// 通知服务相关错误码 (4000-4999)
	NOTIFIER_NOT_FOUND      = 4001 // 通知服务不存在
	NOTIFIER_ALREADY_EXISTS = 4002 // 通知服务已存在
	NOTIFIER_CONFIG_ERROR   = 4003 // 通知服务配置错误
	NOTIFIER_TEST_FAILED    = 4004 // 通知服务测试失败
	NOTIFIER_IN_USE         = 4005 // 通知服务正在使用中

	// 通知发送相关错误码 (5000-5999)
	NOTIFICATION_SEND_FAILED = 5001 // 通知发送失败

	// 插件相关错误码 (6000-6999)
	PLUGIN_NOT_FOUND     = 6001 // 插件不存在
	PLUGIN_DISABLED      = 6002 // 插件未启用
	PLUGIN_CONFIG_ERROR  = 6003 // 插件配置错误
	PLUGIN_LOAD_FAILED   = 6004 // 插件加载失败
	PLUGIN_PROCESS_ERROR = 6005 // 插件处理错误
	PLUGIN_IN_USE        = 6006 // 插件正在使用中

	// 系统错误码 (9000-9999)
	SYSTEM_ERROR        = 9001  // 系统错误
	SERVER_ERROR        = 9002  // 服务器错误
	CONFIG_ERROR        = 9003  // 配置错误
	HEALTH_CHECK_FAILED = 9004  // 健康检查失败
	DEFAULT_ERROR       = 99999 // 默认错误
)
