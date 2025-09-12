package plugin

import (
	"context"
	"emby-plugin/internal/log"
	"emby-plugin/internal/models"
	"emby-plugin/internal/pluginsdk"
	util "emby-plugin/utils"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type EmbyPlugin struct{}

func (p *EmbyPlugin) ID() string      { return "emby" }
func (p *EmbyPlugin) Name() string    { return "Emby Plugin" }
func (p *EmbyPlugin) Version() string { return "1.0.0" }
func (p *EmbyPlugin) Desc() string    { return "解析 Emby Webhook 事件并生成标准输出" }

func (p *EmbyPlugin) DefaultSettings() map[string]any {
	settting := models.Settings{
		IsShowTime:     true,
		IsShowUser:     true,
		IsShowDevice:   true,
		IsShowIP:       true,
		IsShowProgress: true,
		IsShowType:     true,
		IsShowYear:     true,
		PreferURLNames: []string{"MovieDb", "IMDb", "Trakt"},
		Targets:        "",
		ImageSource:    "remote",
		EmbyBaseURL:    "https://127.0.0.1:8096",
		EmbyAPIKey:     "c4c8008b58114d48b57cf442ab1bf301",
		LinkSource:     "remote",
	}
	res := map[string]any{}
	if err := mapstructure.Decode(settting, &res); err != nil {
		log.Logger.Error("设置解码失败", "error", err)
		return map[string]any{}
	}
	return res
}

func (p *EmbyPlugin) Process(ctx context.Context, input map[string]any, settings map[string]any) (*pluginsdk.Output, error) {
	log.Logger.Info("处理Emby事件", "input", input, "settings", settings)
	isNotify := true

	var evt models.EmbyEvent
	decCfg := &mapstructure.DecoderConfig{Result: &evt, TagName: "mapstructure"}
	decoder, err := mapstructure.NewDecoder(decCfg)
	if err != nil {
		return nil, fmt.Errorf("创建解码器失败: %w", err)
	}
	if err := decoder.Decode(input); err != nil {
		return nil, fmt.Errorf("输入解码失败: %w", err)
	}

	var cfg models.Settings
	if err := mapstructure.Decode(settings, &cfg); err != nil {
		return nil, fmt.Errorf("设置解码失败: %w", err)
	}

	notifyEmbyUsers := cfg.NotifyEmbyUsers
	if notifyEmbyUsers != "" && evt.User != nil && strings.TrimSpace(evt.User.Name) != "" {
		users := strings.Split(notifyEmbyUsers, ",")
		if len(users) > 0 {
			isHave := false
			for _, user := range users {
				if user == evt.User.Name {
					isHave = true
					break
				}
			}
			isNotify = isHave
		}
	}
	title := p.buildTitle(evt)
	content := p.buildContent(evt, cfg)
	image := p.buildImage(evt, cfg)
	url := p.buildURL(evt, cfg)
	targets := util.ParseTargets(cfg.Targets)
	output := &pluginsdk.Output{IsNotify: isNotify, Title: title, Content: content, Image: image, URL: url, Targets: targets, Meta: &pluginsdk.MetaData{Req: input, PluginID: p.ID(), ProcessedAt: time.Now().Format(time.RFC3339)}}
	log.Logger.Info("处理Emby事件完成", "output", output)
	return output, nil
}
