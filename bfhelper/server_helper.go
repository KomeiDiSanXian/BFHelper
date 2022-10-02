// 待完善...
package bfhelper
/*
import (
	"regexp"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	// 自动踢出，暂时仅允许群管理使用
	engine.OnPrefixGroup([]string{".自动踢出", "自动踢出", "autokick", ".autokick"}, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			//空格分离各个参数
			str := strings.Split(ctx.State["args"].(string), " ")
			for i := range str {
				limit := regexp.MustCompile(`(rank|kd)[<>=]{1,2}\d+`).FindString(str[i])
			}
		})
}
*/