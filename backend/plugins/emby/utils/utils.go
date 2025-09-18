package utils

import (
	"fmt"
	"strings"
)

// ExternalURL 表示 Emby 中 ExternalUrls 的元素
type ExternalURL struct {
	Name string `mapstructure:"Name" json:"Name"`
	URL  string `mapstructure:"Url" json:"Url"`
}

// PickExternalURL 按优先顺序从 ExternalUrls 中挑选第一个可用链接
func PickExternalURL(urls []ExternalURL, prefer []string) string {
	if len(urls) == 0 {
		return ""
	}
	nameToURL := map[string]string{}
	for _, u := range urls {
		nameToURL[u.Name] = u.URL
	}
	for _, name := range prefer {
		if v, ok := nameToURL[name]; ok && strings.TrimSpace(v) != "" {
			return v
		}
	}
	for _, u := range urls {
		if strings.TrimSpace(u.URL) != "" {
			return u.URL
		}
	}
	return ""
}

// ParseTargets 将 targets 设置解析为字符串数组
func ParseTargets(v any) []string {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		res := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				res = append(res, p)
			}
		}
		return res
	case []string:
		return t
	case []any:
		res := make([]string, 0, len(t))
		for _, x := range t {
			if s, ok := x.(string); ok && strings.TrimSpace(s) != "" {
				res = append(res, strings.TrimSpace(s))
			}
		}
		return res
	default:
		return nil
	}
}

// GetBool 从 map 中安全读取布尔值，未命中返回默认值
func GetBool(m map[string]any, key string, def bool) bool {
	v, ok := m[key]
	if !ok {
		return def
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return def
}

// GetStringSlice 从 map 中读取字符串数组，支持 []string、[]any、逗号分隔字符串
func GetStringSlice(m map[string]any, key string, def []string) []string {
	v, ok := m[key]
	if !ok {
		return def
	}
	if ss, ok := v.([]string); ok {
		return ss
	}
	if aa, ok := v.([]any); ok {
		res := make([]string, 0, len(aa))
		for _, x := range aa {
			if s, ok := x.(string); ok {
				res = append(res, s)
			}
		}
		if len(res) > 0 {
			return res
		}
	}
	if s, ok := v.(string); ok {
		parts := strings.Split(s, ",")
		res := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				res = append(res, p)
			}
		}
		if len(res) > 0 {
			return res
		}
	}
	return def
}

// GetString 从 map 中读取字符串，若不存在返回空字符串
func GetString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

// SafeName 修剪并为缺省名称提供占位
func SafeName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "(未知)"
	}
	return s
}

// SeasonShort 将诸如 "第 1 季" 转换为 "S1"
func SeasonShort(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	var num string
	for _, r := range name {
		if r >= '0' && r <= '9' {
			num += string(r)
		}
	}
	if num == "" {
		return name
	}
	return "S" + num
}

// EpisodeShort 将集序号转换为 "EpN"
func EpisodeShort(index int) string {
	if index <= 0 {
		return ""
	}
	return fmt.Sprintf("Ep%d", index)
}
