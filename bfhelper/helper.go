package bfhelper

import (
	"errors"
	"sync"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"gorm.io/gorm"

	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/record"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 读写锁
var rmu sync.RWMutex

// 引擎注册
var engine = control.Register("战地", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Help:              "",
	PrivateDataFolder: "battlefield",
})

// 返回插件数据目录
func GetDataFolder() string {
	return engine.DataFolder()
}

//权限设置
/*
func permission(ctx *zero.Ctx) bool {
	return ctx.Event.Sender.Role != "member"
	//TODO: 检测该qq是否在权限组中
	/*
		return db.IsExist(grp, ctx.Event.Sender.ID)

}
*/
func init() {
	//初始化数据库
	bf1model.InitDB(engine.DataFolder()+"player.db", &bf1model.Player{})
	//bf1model.InitDB(engine.DataFolder()+"server.db", &bf1model.Admin{}, &bf1model.Server{}, &bf1model.Group{})
	//查询在线玩家数
	engine.OnFullMatchGroup([]string{".bf1stats", "战地1人数", "bf1人数"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("少女折寿中...")
			data, err := rsp.ReturnJson("https://api.s-wg.net/ServersCollection/getStatus", "GET", nil)
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
	engine.OnRegex(`^[\.\/。] *绑定 *(.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			id := ctx.State["regex_matched"].([]string)[1]
			gdb, err := bf1model.Open(engine.DataFolder() + "player.db")
			if err != nil {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败，打开数据库时出错！"))
				return
			}
			db := (*bf1model.PlayerDB)(gdb)
			//先绑定再查询pid和是否实锤
			//检查是否已经绑定
			if data, err := db.FindByQid(ctx.Event.UserID); errors.Is(err, gorm.ErrRecordNotFound) {
				//未绑定...
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("正在绑定id为 ", id))
				rmu.Lock()
				err := db.Create(bf1model.Player{
					Qid:         ctx.Event.UserID,
					DisplayName: id,
				})
				rmu.Unlock()
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
				rmu.Lock()
				err := db.Update(bf1model.Player{
					Qid:         ctx.Event.UserID,
					DisplayName: id,
				})
				rmu.Unlock()
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
				pid, _ = GetPersonalID(id) //TODO:err写入日志
				wg.Done()
			}()
			wg.Wait()
			rmu.Lock()
			db.Update(bf1model.Player{
				PersonalID: pid,
				Qid:        ctx.Event.UserID,
				IsHack:     hack,
			})
			rmu.Unlock()
		})
	// bf1个人战绩
	engine.OnRegex(`^[\.\/。] *1?战绩 *(.*)$`).SetBlock(true).
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
	//所有武器，只展示前五个，修改 RequestWeapon 函数可以展示多个
	engine.OnRegex(`^[\.\/。] *1?武器 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.ALL) })
	//半自动
	engine.OnRegex(`^[\.\/。] *1?半自动 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Semi) })
	//冲锋枪
	engine.OnRegex(`^[\.\/。] *1?冲锋枪 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.SMG) })
	//轻机枪
	engine.OnRegex(`^[\.\/。] *1?轻?机枪 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.LMG) })
	//步枪
	engine.OnRegex(`^[\.\/。] *1?步枪 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Bolt) })
	//霰弹枪
	engine.OnRegex(`^[\.\/。] *1?[霰散]弹枪 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Shotgun) })
	//手枪
	engine.OnRegex(`^[\.\/。] *1?[手配佩]枪 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Sidearm) })
	//近战武器
	engine.OnRegex(`^[\.\/。] *1?近战 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Melee) })
	//手榴弹
	engine.OnRegex(`^[\.\/。] *1?手[榴雷]弹? *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Grenade) })
	//驾驶员
	engine.OnRegex(`^[\.\/。] *1?驾驶员 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Dirver) })
	//配备
	engine.OnRegex(`^[\.\/。] *1?[配装]备 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Gadget) })
	//精英兵
	engine.OnRegex(`^[\.\/。] *1?精英兵? *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) { RequestWeapon(ctx, bf1record.Elite) })
	//Kick 踢出玩家
	/*
		engine.OnRegex(`^[\.\/。] *kick\s*(.*)\s(.*)$`, permission).SetBlock(true).
			Handle(func(ctx *zero.Ctx) {
				id := ctx.State["regex_matched"].([]string)[1]
				reason := ctx.State["regex_matched"].([]string)[2]
				//reason为空的情况
				if id == "" {
					id = reason
					reason = "Adimin kicks!"
				}
				reason = fmt.Sprintf("%s%s", "RemiliaBot: ", reason)
				//检查reason长度
				if len(reason) > 32 {
					ctx.SendChain(message.At(ctx.Event.UserID), message.Text("理由过长！请重新填写！"))
					return
				}
				//translation
				reason = S2tw(reason)
				//.......
			})
	*/
}
