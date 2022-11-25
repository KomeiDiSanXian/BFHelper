package bfhelper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"gorm.io/gorm"

	api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/record"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 引擎注册
var engine = control.Register("战地", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault: false,
	Help: "battlefield\n" +
		"<-----以下是玩家查询----->\n" +
		"- .武器 [武器类型] [id]\t不填武器武器类型默认查询全部\n" +
		"- .载具 [id]\n" +
		"<-----以下是服务器管理----->\n" +
		"- .绑服 gameid1 gameid2...\t仅服主权限可以绑定服务器\n" +
		"- .别名 gameid 别名\t服主权限，需要设置别名，否则部分功能无法使用，注意不要设置相同的别名\n" +
		"- .添加管理 qq1 qq2 qq3...\t仅服主权限可以添加群机器人管理\n" +
		"- .kick player [原因]\t在已绑定的服务器中踢出玩家\n" +
		"- .ban 别名 player\t在别名为此的服务器封禁此玩家\n" +
		"- .unban 别名 player\t在别名为此的服务器解封此玩家\n" +
		"- .自动踢出 -s 别名 -r 等级 -p ping值 -kd kd大小 -kpm kpm大小\t-s 别名必填，踢出大于你所填值的玩家\n" +
		"- 关闭自动踢出\t注意不要加点\n" +
		"<-----以下是更多功能----->\n" +
		"- .bf1stats\t查询亚服相关信息（来自水神的api）\n" +
		"- .交换\t查询本周战地一武器皮肤\n" +
		"- .行动\t查询战地一行动箱子\n" +
		"- .战绩 [id]\t查询生涯的战绩\n" +
		"- .最近 [id]\t查询最近的战绩\n" +
		"- .绑定 id\t进行账号绑定，会检测绑定id是否被实锤",
	PrivateDataFolder: "battlefield",
})

func init() {
	//初始化数据库
	bf1model.InitDB(engine.DataFolder()+"player.db", &bf1model.Player{})
	//查询在线玩家数
	engine.OnFullMatchGroup([]string{".bf1stats", "战地1人数", "bf1人数"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("少女折寿中...")
			data, err := api.ReturnJson("https://api.s-wg.net/ServersCollection/getStatus", "GET", nil)
			if err != nil {
				ctx.Send("ERROR:" + err.Error())
				return
			}
			ctx.SendChain(
				message.At(ctx.Event.UserID),
				message.Text(
					"\n更新时间：", gjson.Get(data, "updateTime").Str, "\n",
					"亚服私服总数：", gjson.Get(data, "quantity"), "\n",
					"在线人数：", gjson.Get(data, "totalPeople"), "[", gjson.Get(data, "totalQueue"), "]", "\n",
					"小模式服：", gjson.Get(data, "mode.smallMode.full"), "/", gjson.Get(data, "mode.smallMode.amount"), "\n",
					"前线：", gjson.Get(data, "mode.frontLine.full"), "/", gjson.Get(data, "mode.frontLine.amount"), "\n",
					"征服：", gjson.Get(data, "mode.conquer.full"), "/", gjson.Get(data, "mode.conquer.amount"), "\n",
					"行动：", gjson.Get(data, "mode.operation.full"), "/", gjson.Get(data, "mode.operation.amount"),
				))
		})
	//Bind QQ绑定ID
	engine.OnPrefixGroup([]string{".绑定", ".bind"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.State["args"].(string)
			//验证id是否有效
			if vld, err := IsValidId(id); vld {
				if err != nil {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				}
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("id无效，请检查id..."))
				return
			}
			gdb, err := bf1model.Open(engine.DataFolder() + "player.db")
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败，打开数据库时出错！"))
				return
			}
			db := (*bf1model.PlayerDB)(gdb)
			defer db.Close()
			//先绑定再查询pid和是否实锤
			//检查是否已经绑定
			if data, err := db.FindByQid(ctx.Event.UserID); errors.Is(err, gorm.ErrRecordNotFound) {
				//未绑定...
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("正在绑定id为 ", id))
				err := db.Create(bf1model.Player{
					Qid:         ctx.Event.UserID,
					DisplayName: id,
				})
				if err != nil {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败，ERR:", err))
					return
				}
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定成功"))
			} else {
				//已绑定，换绑...
				if data.DisplayName == id { //同id..
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("笨蛋！你现在绑的就是这个id"))
					return
				} else {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("将原绑定id为 ", data.DisplayName, " 改绑为 ", id))
				}
				err := db.Update(bf1model.Player{
					Qid:         ctx.Event.UserID,
					DisplayName: id,
				})
				if err != nil {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败，ERR:", err))
					return
				}
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定成功！"))
			}
			hack := false
			pid := ""
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				hack = IsGetBan(id)
				if hack {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("你刚才绑定id: ", id, " 已被联ban实锤！"))
				}
				wg.Done()
			}()
			go func() {
				pid, _ = api.GetPersonalID(id) //TODO:err写入日志
				wg.Done()
			}()
			wg.Wait()
			db.Update(bf1model.Player{
				PersonalID: pid,
				Qid:        ctx.Event.UserID,
				IsHack:     hack,
			})
		})
	// bf1个人战绩
	engine.OnRegex(`\. *1?战绩 *(.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.State["regex_matched"].([]string)[1]
			ctx.Send("少女折寿中...")
			id, err := ReturnBindID(ctx, id)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				return
			}
			stat, err := bf1record.GetStats(id)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("获取失败：", err))
				return
			}
			if stat.Rank == "" {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("获取到的部分数据为空，请检查id是否有效"))
				return
			}
			txt := "id：" + id + "\n" +
				"等级：" + stat.Rank + "\n" +
				"技巧值：" + stat.Skill + "\n" +
				"游玩时长：" + stat.TimePlayed + "\n" +
				"总kd：" + stat.TotalKD + "(" + stat.Kills + "/" + stat.Deaths + ")" + "\n" +
				"总kpm：" + stat.KPM + "\n" +
				"准度：" + stat.Accuracy + "\n" +
				"爆头数：" + stat.Headshots + "\n" +
				"胜率：" + stat.WinPercent + "(" + stat.Wins + "/" + stat.Losses + ")" + "\n" +
				"场均击杀：" + stat.KillsPerGame + "\n" +
				"步战kd：" + stat.InfantryKD + "\n" +
				"步战击杀：" + stat.InfantryKills + "\n" +
				"步战kpm：" + stat.InfantryKPM + "\n" +
				"载具击杀：" + stat.VehicleKills + "\n" +
				"载具kpm：" + stat.VehicleKPM + "\n" +
				"近战击杀：" + stat.DogtagsTaken + "\n" +
				"最高连杀：" + stat.HighestKillStreak + "\n" +
				"最远爆头：" + stat.LongestHeadshot + "\n" +
				"MVP数：" + stat.MVP + "\n" +
				"作为神医拉起了 " + stat.Revives + " 人" + "\n" +
				"开棺材车创死了 " + stat.CarriersKills + " 人"
			Txt2Img(ctx, txt)
		})
	//武器查询，只展示前五个
	engine.OnRegex(`^\. *1?武器 *(.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			str := strings.Split(ctx.State["regex_matched"].([]string)[1], " ")
			id := ""
			if str[0] == "" {
				RequestWeapon(ctx, id, bf1record.ALL)
				return
			}
			//检查str长度
			if len(str) > 1 {
				id = str[1]
			}
			switch str[0] {
			case "半自动", "semi":
				RequestWeapon(ctx, id, bf1record.Semi)
			case "冲锋枪", "冲锋":
				RequestWeapon(ctx, id, bf1record.SMG)
			case "轻机枪", "机枪":
				RequestWeapon(ctx, id, bf1record.LMG)
			case "步枪", "狙击枪", "狙击":
				RequestWeapon(ctx, id, bf1record.Bolt)
			case "霰弹枪", "散弹枪", "霰弹", "散弹":
				RequestWeapon(ctx, id, bf1record.Shotgun)
			case "配枪", "手枪", "副手":
				RequestWeapon(ctx, id, bf1record.Sidearm)
			case "近战", "刀":
				RequestWeapon(ctx, id, bf1record.Melee)
			case "手榴弹", "手雷", "雷":
				RequestWeapon(ctx, id, bf1record.Grenade)
			case "驾驶员", "坦克兵", "载具":
				RequestWeapon(ctx, id, bf1record.Dirver)
			case "配备", "装备":
				RequestWeapon(ctx, id, bf1record.Gadget)
			case "精英", "精英兵":
				RequestWeapon(ctx, id, bf1record.Elite)
			default:
				if regexp.MustCompile(`\w+`).MatchString(str[0]) {
					id = str[0]
					RequestWeapon(ctx, id, bf1record.ALL)
				}
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("请检查输入格式是否有误...")))
			}
		})
	//最近战绩
	engine.OnRegex(`^\. *1?最近 *(.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.State["regex_matched"].([]string)[1]
			ctx.Send("少女折寿中...")
			id, err := ReturnBindID(ctx, id)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				return
			}
			recent, err := GetBF1Recent(id)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				return
			}
			msg := "id：" + id + "\n"
			for i := range *recent {
				msg += "服务器：" + (*recent)[i].Server[:24] + "\n"
				msg += "地图：" + (*recent)[i].Map + "   (" + (*recent)[i].Mode + ")\n"
				msg += "kd：" + strconv.FormatFloat((*recent)[i].Kd, 'f', -1, 64) + "\n"
				msg += "kpm：" + strconv.FormatFloat((*recent)[i].Kpm, 'f', -1, 64) + "\n"
				msg += "游玩时长：" + strconv.FormatFloat(float64((*recent)[i].Time/60), 'f', -1, 64) + "分钟"
				msg += "\n---------------\n"
			}
			Txt2Img(ctx, msg)
		})
	//获取所有种类的载具信息
	engine.OnRegex(`^\. *1?载具 *(.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.State["regex_matched"].([]string)[1]
			ctx.Send("少女折寿中...")
			pid, id, err := ID2PID(ctx.Event.UserID, id)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				return
			}
			car, err := bf1record.GetVehicles(pid)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				return
			}
			msg := "id：" + id + "\n"
			for i := range *car {
				msg += "------------\n"
				msg += (*car)[i].Name + "\n"
				msg += fmt.Sprintf("%s%6.0f\t", "击杀数：", (*car)[i].Kills)
				msg += "kpm：" + (*car)[i].KPM + "\n"
				msg += fmt.Sprintf("%s%6.0f\t", "击毁数：", (*car)[i].Destroyed)
				msg += "游玩时间：" + (*car)[i].Time + " 小时\n"
			}
			Txt2Img(ctx, msg)
		})
	//交换查询
	engine.OnFullMatchGroup([]string{".交换", ".exchange"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			exchange, err := api.GetExchange()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			var msg string
			for i, v := range exchange {
				msg += i + "：\n"
				for _, skin := range v {
					msg += "\t" + skin + "\n"
				}
			}
			Txt2Img(ctx, msg)
		})
	//行动包查询
	engine.OnFullMatchGroup([]string{".行动", ".行动包", ".pack"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			pack, err := api.GetCampaignPacks()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			var msg string
			msg += "行动名：" + pack.Name + "\n"
			msg += "剩余时间：" + fmt.Sprintf("%.2f", float64(pack.RemainTime)/60/24) + " 天\n"
			msg += "箱子重置时间：" + fmt.Sprintf("%.2f", float64(pack.ResetTime)/60) + " 小时\n"
			msg += "行动地图：" + pack.Op1Name + " 与 " + pack.Op2Name + "\n"
			msg += "行动简介：" + pack.Desc
			Txt2Img(ctx, msg)
		})
}
