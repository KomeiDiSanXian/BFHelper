// Package bfhelper 战地玩家查询
package bfhelper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/jinzhu/gorm"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	bf1record "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/record"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/setting"
)

// 引擎注册
var engine = control.Register("战地", &ctrl.Options[*zero.Ctx]{
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

// EngineFile 返回 engine.DataFolder()
func EngineFile() string {
	return engine.DataFolder()
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	if err := setting.ReadSection("Account", &global.Account.LoginedUser); err != nil {
		return err
	}
	return setting.ReadSection("SakuraKooi", &global.SakuraAPI)
}

var dbname = engine.DataFolder() + "battlefield.db"

func init() {
	pluginInitSuccess := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		// 初始化数据库
		if err := bf1model.Init(dbname); err != nil {
			ctx.SendChain(message.Text("ERROR: 数据库初始化失败, 请联系机器人管理员重启"))
			return false
		}
		// 读取配置文件
		if err := setupSetting(); err != nil {
			ctx.SendChain(message.Text("ERROR: 读取插件配置失败, 请联系机器人管理员重启"))
			return false
		}
		// 刷新Session
		_ = api.Login(global.Account.LoginedUser.Username, global.Account.LoginedUser.Password, true)
		return true
	})

	// Bind QQ绑定ID
	engine.OnPrefixGroup([]string{".绑定", ".bind"}, pluginInitSuccess).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.State["args"].(string)
			// 验证id是否有效
			if vld, err := IsValidID(id); vld {
				if err != nil {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
				}
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("id无效，请检查id..."))
				return
			}
			db, err := bf1model.Open(dbname)
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败，打开数据库时出错！"))
				return
			}
			defer db.Close()
			// 先绑定再查询pid和是否实锤
			// 检查是否已经绑定
			playerRepo := bf1model.NewPlayerRepository(db)
			if data, err := playerRepo.GetByQID(ctx.Event.UserID); errors.Is(err, gorm.ErrRecordNotFound) {
				// 未绑定...
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("正在绑定id为 ", id))
				err := playerRepo.Create(&bf1model.Player{
					Qid:         ctx.Event.UserID,
					DisplayName: id,
				})
				if err != nil {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败，ERR:", err))
					return
				}
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定成功"))
			} else {
				// 已绑定，换绑...
				if data.DisplayName == id { // 同id..
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("笨蛋！你现在绑的就是这个id"))
					return
				}
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("将原绑定id为 ", data.DisplayName, " 改绑为 ", id))
				err := playerRepo.Update(&bf1model.Player{
					Qid:         ctx.Event.UserID,
					DisplayName: id,
				})
				if err != nil {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败，ERR:", err))
					return
				}
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定成功！"))
			}
			pid, _ := api.GetPersonalID(id) //TODO:err写入日志
			_ = playerRepo.Update(&bf1model.Player{
				PersonalID: pid,
				Qid:        ctx.Event.UserID,
			})
		})
	// bf1个人战绩
	engine.OnRegex(`\. *1?战绩 *(.*)$`, pluginInitSuccess).SetBlock(true).
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
	// 武器查询，只展示前五个
	engine.OnRegex(`^\. *1?武器 *(.*)$`, pluginInitSuccess).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			str := strings.Split(ctx.State["regex_matched"].([]string)[1], " ")
			id := ""
			if str[0] == "" {
				RequestWeapon(ctx, id, bf1record.ALL)
				return
			}
			// 检查str长度
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
			}
		})
	// 最近战绩
	engine.OnRegex(`^\. *1?最近 *(.*)$`, pluginInitSuccess).SetBlock(true).
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
	// 获取所有种类的载具信息
	engine.OnRegex(`^\. *1?载具 *(.*)$`, pluginInitSuccess).SetBlock(true).
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
	// 交换查询
	engine.OnFullMatchGroup([]string{".交换", ".exchange"}, pluginInitSuccess).SetBlock(true).
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
	// 行动包查询
	engine.OnFullMatchGroup([]string{".行动", ".行动包", ".pack"}, pluginInitSuccess).SetBlock(true).
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
