package plugin

import (
	"emby-plugin/internal/emby"
	"emby-plugin/internal/log"
	"emby-plugin/internal/models"
	"fmt"
	"strings"
)

func (p *EmbyPlugin) buildImage(evt models.EmbyEvent, settings models.Settings) string {
	if !(strings.HasPrefix(evt.Event, "playback") || strings.HasPrefix(evt.Event, "library")) {
		return ""
	}
	id := strings.TrimSpace(evt.Item.ID)
	if id == "" {
		return ""
	}
	imageSource := strings.ToLower(strings.TrimSpace(settings.ImageSource)) // nas | remote
	if imageSource == "remote" {
		emby := emby.NewEmby(settings.EmbyBaseURL, settings.EmbyAPIKey)
		log.Logger.Info("获取远程图片", "id", id)
		image, err := emby.FetchEmbyRemoteImageURL(id, "Backdrop")
		if err == nil && image != "" {
			return image
		} else {
			log.Logger.Error("获取远程图片失败", "error", err)
		}
	}
	base := strings.TrimRight(settings.EmbyBaseURL, "/")
	if base == "" {
		return ""
	}
	if len(evt.Item.BackdropImageTags) > 0 {
		tag := strings.TrimSpace(evt.Item.BackdropImageTags[0])
		if tag != "" {
			return fmt.Sprintf("%s/Items/%s/Images/Backdrop/0?tag=%s&quality=90", base, id, tag)
		}
	}
	if evt.Item.ImageTags != nil {
		if primary, ok := evt.Item.ImageTags["Primary"]; ok && strings.TrimSpace(primary) != "" {
			return fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s&quality=90", base, id, primary)
		}
	}
	return ""
}
