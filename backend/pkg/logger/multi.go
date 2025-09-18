package logger

import (
	"context"
	"log/slog"
)

type multiHandler struct{ hs []slog.Handler }

func (m multiHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	for _, h := range m.hs {
		if h.Enabled(ctx, lvl) {
			return true
		}
	}
	return false
}
func (m multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.hs {
		_ = h.Handle(ctx, r) // 忽略单个 handler 出错
	}
	return nil
}
func (m multiHandler) WithAttrs(a []slog.Attr) slog.Handler {
	nh := make([]slog.Handler, len(m.hs))
	for i, h := range m.hs {
		nh[i] = h.WithAttrs(a)
	}
	return multiHandler{nh}
}
func (m multiHandler) WithGroup(g string) slog.Handler {
	nh := make([]slog.Handler, len(m.hs))
	for i, h := range m.hs {
		nh[i] = h.WithGroup(g)
	}
	return multiHandler{nh}
}
