// Package engine 插件注册
package engine

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// Engine 引擎注册
var Engine = control.Register("战地", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault: false,
	Brief:            "战地相关查询功能",
	Help: "battlefield\n" +
		"<-----以下是玩家查询----->\n" +
		"- .武器 [武器类型] [id]\t不填武器武器类型默认查询全部\n" +
		"- .载具 [id]\n" +
		"<-----以下是更多功能----->\n" +
		"- .交换\t查询本周战地一武器皮肤\n" +
		"- .行动\t查询战地一行动箱子\n" +
		"- .战绩 [id]\t查询生涯的战绩\n" +
		"- .最近 [id]\t查询最近的战绩\n" +
		"- .绑定 id\t进行账号绑定，会检测绑定id是否被实锤",
	PrivateDataFolder: "battlefield",
})