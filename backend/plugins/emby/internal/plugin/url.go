package plugin

import (
	"emby-plugin/internal/models"
	util "emby-plugin/utils"
	"fmt"
	"strings"
)

func (p *EmbyPlugin) buildURL(evt models.EmbyEvent, settings models.Settings) string {
	switch strings.ToLower(strings.TrimSpace(settings.LinkSource)) {
	case "emby":
		base := strings.TrimRight(settings.EmbyBaseURL, "/")
		if base == "" || strings.TrimSpace(evt.Item.ID) == "" {
			return ""
		}
		return fmt.Sprintf("%s/web/index.html#!/item?id=%s", base, evt.Item.ID)
	default:
		prefer := settings.PreferURLNames
		if len(prefer) == 0 {
			prefer = []string{"MovieDb", "IMDb", "Trakt"}
		}
		return util.PickExternalURL(evt.Item.ExternalUrls, prefer)
	}
}
