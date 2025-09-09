package server

import (
	"net/http"

	"github.com/jianxcao/notify/internal/logger"

	"github.com/gin-gonic/gin"
)

// setupLogRoutes 设置日志流相关路由
func (s *HTTPServer) setupLogRoutes(api *gin.RouterGroup) {
	api.GET("/logs/stream", s.handleLogStream)
}

// handleLogStream 处理日志 SSE 流
func (s *HTTPServer) handleLogStream(c *gin.Context) {
	// 设置SSE相关头部
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 订阅日志
	ch := logger.Subscribe()
	defer logger.Unsubscribe(ch)

	// 将请求上下文用于取消
	ctx := c.Request.Context()

	// 简单心跳，防止某些代理超时
	c.Writer.Flush()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// 发送事件
			c.SSEvent("message", msg)
			// 手动刷新
			if f, ok := c.Writer.(http.Flusher); ok {
				f.Flush()
			}
		case <-ctx.Done():
			return
		}
	}
}
