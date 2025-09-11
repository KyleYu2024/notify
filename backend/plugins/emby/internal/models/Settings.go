package models

// Settings 定义插件的可配置项，通过 mapstructure 从 map[string]any 解码
type Settings struct {
	PreferURLNames []string `mapstructure:"prefer_url_names" json:"prefer_url_names"`
	Targets        string   `mapstructure:"targets" json:"targets"`
	ImageSource    string   `mapstructure:"image_source" json:"image_source"`
	EmbyBaseURL    string   `mapstructure:"emby_base_url" json:"emby_base_url"`
	EmbyAPIKey     string   `mapstructure:"emby_api_key" json:"emby_api_key"`
	LinkSource     string   `mapstructure:"link_source" json:"link_source"`
	IsShowTime     bool     `mapstructure:"is_show_time" json:"is_show_time"`
	IsShowUser     bool     `mapstructure:"is_show_user" json:"is_show_user"`
	IsShowDevice   bool     `mapstructure:"is_show_device" json:"is_show_device"`
	IsShowIP       bool     `mapstructure:"is_show_ip" json:"is_show_ip"`
	IsShowProgress bool     `mapstructure:"is_show_progress" json:"is_show_progress"`
	IsShowType     bool     `mapstructure:"is_show_type" json:"is_show_type"`
	IsShowYear     bool     `mapstructure:"is_show_year" json:"is_show_year"`
}
