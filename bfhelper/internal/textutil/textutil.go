// Package textutil 用于处理文字
package textutil

import (
	"strings"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
)

// Traditionalize 简体转繁体
func Traditionalize(text string) string {
	result := ""
	for _, v := range text {
		r, ok := global.Dictionary[string(v)]
		if ok {
			result += r
			continue
		}
		result += string(v)
	}
	return result
}

// CleanPersonalID 检查是否为pid (有 # 就判断为pid)
//
// 如果是就删去 # 并返回 pid
//
// 为否就返回输入
func CleanPersonalID(input string) (cleaned string, hasHash bool) {
	if strings.Contains(input, "#") {
		cleaned = strings.Trim(input, "#")
		return cleaned, true
	}
	return input, false
}

