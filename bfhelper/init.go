package bfhelper

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

//插件引擎注册
var engine = control.Register("战地", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	PrivateDataFolder: "Battlefield",
})

func init() {
	engine.OnFullMatchGroup([]string{".bf1stats", "战地1人数", "bf1人数"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
		})
}
