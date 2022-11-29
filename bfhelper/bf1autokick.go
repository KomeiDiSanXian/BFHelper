package bfhelper

import (
	"fmt"
	"strconv"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/record"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/rsp"
	"github.com/fumiama/cron"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// flag
type autokick struct {
	Server string  `flag:"s"`
	Rank   int     `flag:"r"`
	Ping   int     `flag:"p"`
	Kd     float64 `flag:"kd"`
	Kpm    float64 `flag:"kpm"`
}

// 白名单pid数组
type whitelists []string

var en = control.Register("bf1自动踢出", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault: false,
	Help: "bf1autokick\n" +
		"- .自动踢出 -s 别名 -r 等级 -p ping值 -kd kd大小 -kpm kpm大小\t-s 别名必填，踢出大于你所填值的玩家\n" +
		"- 关闭自动踢出\t注意不要加点\n",
}).ApplySingle(single.New(
	single.WithKeyFn(func(ctx *zero.Ctx) int64 { return ctx.Event.GroupID }),
	single.WithPostFn[int64](func(ctx *zero.Ctx) {
		ctx.Send(
			message.ReplyWithMessage(ctx.Event.MessageID,
				message.Text("自动踢出已开启..."),
			),
		)
	}),
))

func init() {
	en.OnShell("自动踢出", autokick{}, zero.OnlyGroup, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			info := ctx.State["flag"].(*autokick)
			// 检查参数
			if info.Server == "" {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：未指定服务器"))
				return
			}
			var msg message.Message
			msg = append(msg, message.Reply(ctx.Event.MessageID), message.Text("将在 ", info.Server, " 踢出"))
			if info.Rank != 0 {
				msg = append(msg, message.Text("等级大于", info.Rank, " "))
			} else {
				info.Rank = 151
			}
			if info.Ping != 0 {
				msg = append(msg, message.Text("ping大于", info.Ping, " "))
			} else {
				info.Ping = 9999
			}
			if info.Kd != 0 {
				msg = append(msg, message.Text("kd大于", info.Kd, " "))
			} else {
				info.Kd = 999
			}
			if info.Kpm != 0 {
				msg = append(msg, message.Text("kpm大于", info.Kpm, " "))
			} else {
				info.Kpm = 999
			}
			msg = append(msg, message.Text("的玩家"))
			ctx.SendChain(msg...)
			db, cl, err := OpenServerDB()
			defer cl()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			s, err := db.GetServer(info.Server, ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			type p struct {
				GameIds []string `json:"gameIds"`
			}
			post := p{
				GameIds: []string{s.Gameid},
			}
			var sym chan bool
			var wl whitelists
			srv := bf1rsp.NewServer(s.Serverid, s.Gameid, s.PGid)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在获取服务器管理员名单..."))
			admins, err := srv.GetAdminspid()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			wl.AddWhitelist(admins...)
			c := cron.New()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已添加服务器管理员到白名单，开始自动踢出"))
			go func() {
				c.AddFunc("@every 60s", func() {
					data, err := bf1api.ReturnJson(bf1api.EASBAPI, "POST", &post)
					if err != nil {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
						return
					}
					players := gjson.Get(data, "data."+s.Gameid+".players")
					if len(players.Array()) < 20 {
						c.Stop()
						sym <- true
						ctx.SendChain(message.Text(info.Server, "中人数不足20，正在关闭自动踢出..."))
						return
					}
					players.ForEach(func(_, value gjson.Result) bool {
						go func(value gjson.Result) {
							pid := strconv.FormatInt(value.Get("pid").Int(), 10)
							if wl.IsInWhitelist(pid) {
								return
							}
							if value.Get("rank").Int() > int64(info.Rank) {
								ctx.SendChain(message.Text("正在踢出", value.Get("display_name"), "：等级过高(", value.Get("rank"), ")"))
								srv.Kick(pid, "Rank limit "+strconv.Itoa(info.Rank))
							}
							if value.Get("ping").Int() > int64(info.Ping) {
								ctx.SendChain(message.Text("正在踢出", value.Get("display_name"), "：ping值过高(", value.Get("ping"), ")"))
								srv.Kick(pid, "Ping limit "+strconv.Itoa(info.Ping))
							}
							kd, kpm, err := bf1record.Get2k(pid)
							if err != nil {
								return
							}
							if kd > info.Kd {
								ctx.SendChain(message.Text("正在踢出", value.Get("display_name"), "：kd过高(", fmt.Sprintf("%.2f", kd), ")"))
								srv.Kick(pid, "Life KD limit "+strconv.FormatFloat(info.Kd, 'f', 2, 32))
							}
							if kpm > info.Kpm {
								ctx.SendChain(message.Text("正在踢出", value.Get("display_name"), "：kpm过高(", kpm, ")"))
								srv.Kick(pid, "Life KPM limit "+strconv.FormatFloat(info.Kpm, 'f', 2, 32))
							}
						}(value)
						return true
					})
				})
				c.Start()
			}()
			next := zero.NewFutureEvent("message", 999, false, zero.OnlyGroup, zero.RegexRule(`^关闭自动踢出$`), zero.CheckGroup(ctx.Event.GroupID), ServerAdminPermission)
			recv, cancle := next.Repeat()
			defer cancle()
			for {
				select {
				case r := <-recv:
					c.Stop()
					ctx.SendChain(message.Reply(r.Event.MessageID), message.Text("正在关闭自动踢出..."))
					ctx.SendChain(message.Text(info.Server, " 的自动踢出已关闭"))
					return
				case <-sym:
					c.Stop()
					ctx.SendChain(message.Text(info.Server, " 的自动踢出已关闭"))
					return
				}
			}
		})
}

func (wl *whitelists) AddWhitelist(pids ...string) {
	*wl = append(*wl, pids...)
}

func (wl *whitelists) IsInWhitelist(pid string) bool {
	for _, v := range *wl {
		if pid == v {
			return true
		}
	}
	return false
}
