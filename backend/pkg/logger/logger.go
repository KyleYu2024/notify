package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jianxcao/notify/backend/pkg/config"
)

var Logger *slog.Logger

// CallerSkipHandler wraps another handler and adjusts the source location
type CallerSkipHandler struct {
	handler slog.Handler
	skip    int
}

func NewCallerSkipHandler(handler slog.Handler, skip int) *CallerSkipHandler {
	return &CallerSkipHandler{
		handler: handler,
		skip:    skip,
	}
}

func (h *CallerSkipHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CallerSkipHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CallerSkipHandler{
		handler: h.handler.WithAttrs(attrs),
		skip:    h.skip,
	}
}

func (h *CallerSkipHandler) WithGroup(name string) slog.Handler {
	return &CallerSkipHandler{
		handler: h.handler.WithGroup(name),
		skip:    h.skip,
	}
}

func (h *CallerSkipHandler) Handle(ctx context.Context, r slog.Record) error {
	// 调整PC（Program Counter）来跳过包装函数
	if h.skip > 0 {
		var pcs [1]uintptr
		runtime.Callers(h.skip+1, pcs[:])
		r.PC = pcs[0]
	}
	return h.handler.Handle(ctx, r)
}

// 获取项目根目录路径（缓存避免重复计算）
var projectRoot string

func getProjectRoot() string {
	if projectRoot == "" {
		// 获取当前工作目录
		if wd, err := os.Getwd(); err == nil {
			projectRoot = wd
		} else {
			// 如果获取工作目录失败，尝试通过可执行文件路径获取
			if ex, err := os.Executable(); err == nil {
				projectRoot = filepath.Dir(ex)
			}
		}
		// 确保路径以分隔符结尾，方便后续处理
		if projectRoot != "" && !strings.HasSuffix(projectRoot, string(filepath.Separator)) {
			projectRoot += string(filepath.Separator)
		}
	}
	return projectRoot
}

// Init 初始化日志系统
func Init() {
	// 设置日志级别
	var level slog.Level
	switch config.EnvCfg.LOG_LEVEL {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 创建处理器选项
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // 启用源代码位置信息
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 只对顶级属性进行处理（不在group内的属性）
			if len(groups) == 0 && a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					// 通过项目根目录动态计算相对路径
					root := getProjectRoot()
					if root != "" && strings.HasPrefix(source.File, root) {
						// 去掉项目根目录路径，保留相对路径
						source.File = strings.TrimPrefix(source.File, root)
					} else {
						// 如果无法获取项目根目录或文件不在项目内，只显示文件名
						source.File = filepath.Base(source.File)
					}
				}
			}
			return a
		},
	}

	// 根据格式选择处理器
	var baseHandler slog.Handler
	if config.EnvCfg.LOG_FORMAT == "text" {
		baseHandler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// 默认使用 JSON 格式
		baseHandler = slog.NewJSONHandler(os.Stdout, opts)
	}

	// 包装handler来处理caller skip
	// skip=4 是因为：runtime.Callers -> Handle -> 便捷方法 -> 实际调用者
	handler := NewCallerSkipHandler(baseHandler, 5)
	webHandler := NewCallerSkipHandler(slog.NewJSONHandler(&ChannelWriter{}, opts), 5)
	mh := multiHandler{[]slog.Handler{
		handler,
		webHandler,
	}}
	// 创建并设置全局 logger
	Logger = slog.New(mh)
	slog.SetDefault(Logger)
}

// 便捷方法 - 现在会正确显示调用者位置
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

// InfoWithContext 带上下文的日志方法
func InfoWithContext(msg string, contextFields map[string]any, args ...any) {
	allArgs := make([]any, 0, len(args)+len(contextFields)*2)

	// 添加上下文字段
	for k, v := range contextFields {
		allArgs = append(allArgs, k, v)
	}

	// 添加其他参数
	allArgs = append(allArgs, args...)

	Logger.Info(msg, allArgs...)
}

func ErrorWithContext(msg string, contextFields map[string]any, args ...any) {
	allArgs := make([]any, 0, len(args)+len(contextFields)*2)

	// 添加上下文字段
	for k, v := range contextFields {
		allArgs = append(allArgs, k, v)
	}

	// 添加其他参数
	allArgs = append(allArgs, args...)

	Logger.Error(msg, allArgs...)
}

// Fatal 记录错误日志并退出程序
func Fatal(msg string, args ...any) {
	Logger.Error(msg, args...)
	os.Exit(1)
}
