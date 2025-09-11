package plugin

import (
	"emby-plugin/internal/models"
	"fmt"
	"strings"
	"time"
)

func (p *EmbyPlugin) buildContent(evt models.EmbyEvent, settings models.Settings) string {
	if !(strings.HasPrefix(evt.Event, "playback") || strings.HasPrefix(evt.Event, "library")) {
		return ""
	}
	var b strings.Builder
	if settings.IsShowTime && strings.TrimSpace(evt.Date) != "" {
		if t, err := time.Parse(time.RFC3339Nano, evt.Date); err == nil {
			b.WriteString("📅 时间: ")
			b.WriteString(t.Local().Format("2006-01-02 15:04:05"))
			b.WriteString("\n")
		}
	}

	if settings.IsShowUser && evt.User != nil && strings.TrimSpace(evt.User.Name) != "" {
		b.WriteString("👤 用户: ")
		b.WriteString(evt.User.Name)
		b.WriteString("\n")
	}
	if evt.Session != nil {
		dev := strings.TrimSpace(evt.Session.DeviceName)
		cli := strings.TrimSpace(evt.Session.Client)
		if settings.IsShowDevice && dev != "" && cli != "" {
			b.WriteString("📱 设备: ")
			b.WriteString(dev)
			b.WriteString(" (")
			b.WriteString(cli)
			b.WriteString(")\n")
		} else if settings.IsShowDevice && dev != "" {
			b.WriteString("📱 设备: ")
			b.WriteString(dev)
			b.WriteString("\n")
		}
		if settings.IsShowIP && strings.TrimSpace(evt.Session.RemoteEndPoint) != "" {
			b.WriteString("🌐 IP: ")
			b.WriteString(evt.Session.RemoteEndPoint)
			b.WriteString("\n")
		}
	}

	if settings.IsShowProgress {
		if evt.TranscodingInfo != nil {
			b.WriteString("⏱️ 播放进度: ")
			b.WriteString(fmt.Sprintf("%.1f%%", evt.TranscodingInfo.CompletionPercentage))
			b.WriteString("\n")
		}
		if evt.PlaybackInfo != nil && evt.Item.RunTimeTicks > 0 && evt.PlaybackInfo.PositionTicks >= 0 {
			pct := float64(evt.PlaybackInfo.PositionTicks) / float64(evt.Item.RunTimeTicks) * 100
			b.WriteString("⏱️ 播放进度: ")
			b.WriteString(fmt.Sprintf("%.1f%%", pct))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(p.buildEpisodeInfo(evt))
	b.WriteString("\n")

	if settings.IsShowType {
		b.WriteString("🎯 类型: ")
		b.WriteString(p.buildTypeInfo(evt))
		b.WriteString("\n")
	}

	if settings.IsShowYear {
		b.WriteString("🗓️ 年份: ")
		b.WriteString(fmt.Sprintf("%d", evt.Item.ProductionYear))
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}
