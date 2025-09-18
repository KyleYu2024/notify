package utils

import (
	"regexp"
	"strings"
)

// 中国大陆手机号码正则表达式
var mobilePhoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// IsMobilePhone 检测字符串是否为有效的中国大陆手机号码
func IsMobilePhone(phone string) bool {
	// 去除可能的空格和连字符
	cleanPhone := strings.ReplaceAll(strings.ReplaceAll(phone, " ", ""), "-", "")
	return mobilePhoneRegex.MatchString(cleanPhone)
}
