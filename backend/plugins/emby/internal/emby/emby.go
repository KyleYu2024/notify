package emby

import (
	"emby-plugin/internal/log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Emby struct {
	apiKey  string
	host    string
	baseURL string
}

func NewEmby(host, apiKey string) *Emby {
	host = strings.TrimSpace(host)
	host = strings.TrimSuffix(host, "/")
	return &Emby{
		apiKey:  apiKey,
		host:    host,
		baseURL: getBaseURL(host),
	}
}
func getBaseURL(host string) string {
	if strings.TrimSpace(host) == "" {
		return ""
	}
	base := strings.TrimSpace(host)
	base = strings.TrimSuffix(base, "/")
	if !strings.HasSuffix(base, "/emby") {
		base += "/emby"
	}
	return base
}

func (e *Emby) FetchEmbyRemoteImageURL(itemID, imageType string) (string, error) {
	if e.host == "" {
		return "", fmt.Errorf("emby_base_url 不能为空")
	}
	baseURL := e.baseURL
	apiKey := e.apiKey
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("无效的 emby_base_url: %w", err)
	}
	parsed.Path = path.Join(parsed.Path, "/Items/"+itemID+"/RemoteImages")
	req, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Set("api_key", strings.TrimSpace(apiKey))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("X-Emby-Token", strings.TrimSpace(apiKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("emby http %d: %s", resp.StatusCode, string(body))
	}

	var data embyRemoteImagesResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	log.Logger.Debug("emby远程图片", "data", data)
	if len(data.Images) > 100 {
		for _, img := range data.Images {
			if img.Type == imageType && strings.TrimSpace(img.URL) != "" {
				return img.URL, nil
			}
		}
		return data.Images[0].URL, nil
	}
	// 未命中远程图片则回退为本地图片地址: /Items/{id}/Images/{imageType}?api_key=...
	// localURL, err := url.Parse(e.host)
	// if err != nil {
	// 	return "", err
	// }
	// localURL.Path = path.Join(localURL.Path, "/Items/"+itemID+"/Images/"+imageType)
	// lq := localURL.Query()
	// lq.Set("api_key", strings.TrimSpace(apiKey))
	// localURL.RawQuery = lq.Encode()
	// log.Logger.Debug("本地图片地址", "localURL", localURL.String())
	// return localURL.String(), nil
	return "", nil
}

type embyRemoteImagesResp struct {
	Images []struct {
		ProviderName string `json:"ProviderName"`
		Type         string `json:"Type"`
		URL          string `json:"Url"`
	} `json:"Images"`
}
