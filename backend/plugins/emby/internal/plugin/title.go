package plugin

import (
	"emby-plugin/internal/models"
	"fmt"
	"strings"
)

func (p *EmbyPlugin) buildTitle(evt models.EmbyEvent) string {
	var sb strings.Builder
	if evt.User != nil && strings.TrimSpace(evt.User.Name) != "" {
		sb.WriteString(evt.User.Name)
		sb.WriteString(" ")
	}
	switch evt.Event {
	case "playback.start":
		sb.WriteString("开始播放 ")
	case "playback.stop", "playback.unpause", "playback.pause":
		sb.WriteString("停止播放 ")
	case "library.new":
		sb.WriteString("新增媒体  ")
	case "library.updated":
		sb.WriteString("更新媒体 ")
	case "library.deleted":
		sb.WriteString("删除媒体 ")
	case "system.notificationtest":
		sb.WriteString("测试媒体 ")
	default:
		sb.WriteString(evt.Title)
		sb.WriteString(" ")
	}
	sb.WriteString(p.buildEpisodeInfo(evt))
	return strings.TrimSpace(sb.String())
}

func (p *EmbyPlugin) buildEpisodeInfo(evt models.EmbyEvent) string {
	var sb strings.Builder
	switch evt.Item.Type {
	case "Episode", "Series", "Season":
		sb.WriteString("📺 剧集:")
		item := evt.Item
		if item.SeriesName != "" && item.IndexNumber > 0 && item.ParentIndexNumber > 0 {
			sb.WriteString(fmt.Sprintf("%s 第%d季 第%d集", item.SeriesName, item.IndexNumber, item.ParentIndexNumber))
		} else if item.SeriesName != "" {
			sb.WriteString(item.SeriesName)
			if item.SeasonName != "" {
				sb.WriteString(" ")
				sb.WriteString(item.SeasonName)
			}
			if item.Name != "" {
				sb.WriteString(" ")
				sb.WriteString(item.Name)
			}
		} else if item.Name != "" {
			sb.WriteString(item.Name)
		}
		sb.WriteString("\n")
	case "Audio":
		sb.WriteString("🎧 音频: ")
		sb.WriteString(evt.Item.Album)
		sb.WriteString("\n")
	case "Movie":
		sb.WriteString("🎬 电影: ")
		sb.WriteString(evt.Item.Name)
		if evt.Item.ProductionYear > 0 {
			sb.WriteString(fmt.Sprintf(" (%d)", evt.Item.ProductionYear))
		}
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func (p *EmbyPlugin) buildTypeInfo(evt models.EmbyEvent) string {
	var sb strings.Builder
	switch evt.Item.Type {
	case "Episode", "Series", "Season":
		sb.WriteString("电视剧")
	case "Movie":
		sb.WriteString("电影")
	case "Audio":
		sb.WriteString("音频")
	}
	return strings.TrimSpace(sb.String())
}
