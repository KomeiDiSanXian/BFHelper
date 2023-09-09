// Package global 插件注册
package global

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
		"- .cb [id] 查询玩家EAC和BFBan案件信息\n" +
		"<-----以下是服务器管理----->\n" +
		"- .k id 将 id 踢出服务器\n" +
		"- .b alias id 在别名为 alias 的服务器封禁 id\n" +
		"- .ub alias id 在别名为 alias 的服务器解封 id\n" +
		"- .bana alias id 在所有已绑定的服务器封禁 id\n" +
		"- .ubana alias id 在所有已绑定的服务器解封 id\n" +
		"- .cm alias [地图id] 在别名为 alias 的服务器切换地图到地图id\n" +
		"- .qm alias 查询别名为 alias 的地图池信息\n" +
		"<-----以下是更多功能----->\n" +
		"- .交换\t查询本周战地一武器皮肤\n" +
		"- .行动\t查询战地一行动箱子\n" +
		"- .战绩 [id]\t查询生涯的战绩\n" +
		"- .最近 [id]\t查询最近的战绩\n" +
		"- .绑定 id\t进行账号绑定",
	PrivateDataFolder: "battlefield",
})
