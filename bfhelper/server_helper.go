// 待完善...
package bfhelper
/*
import (
	"regexp"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
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
	engine.OnPrefix(".查水表", zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.State["args"].(string)
			// 检查id有效性
			if vld, err := IsValidId(id); vld {
				if err != nil {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				}
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("id无效，请检查id..."))
				return
			}
		})
}
*/