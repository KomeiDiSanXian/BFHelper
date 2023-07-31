// Package textutil 用于处理文字
package textutil

import "github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"

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
