// Package renderer 文字转图片并发送
package renderer

import (
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// Txt2Img 文字转图片并发送
func Txt2Img(ctx *zero.Ctx, txt string) {
	data, err := text.RenderToBase64(txt, text.FontFile, 400, 20)
	if err != nil {
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("将文字转换成图片时发生错误")))
	}
	if id := ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Image("base64://"+helper.BytesToString(data)))); id.ID() == 0 {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERROR:可能被风控了"))
	}
}
