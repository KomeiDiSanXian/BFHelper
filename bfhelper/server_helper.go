// 待完善...
package bfhelper

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	bf1rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/rsp"
	"github.com/fumiama/cron"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
)

// flag
type autokick struct {
	Server string `flag:"s"`
	Rank   int    `flag:"r"`
	Ping   int    `flag:"p"`
}

func init() {
	bf1model.InitDB(engine.DataFolder()+"server.db", &bf1model.Group{}, &bf1model.Server{}, &bf1model.Admin{})

	engine.OnPrefix(".init", zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// args[0] groupid
			// args[1] ownerid
			// args[2] gameid
			args := strings.Split(ctx.State["args"].(string), " ")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在新建..."))
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			grpid, _ := strconv.ParseInt(args[0], 10, 64)
			ownerid, _ := strconv.ParseInt(args[1], 10, 64)
			err = db.Create(grpid, ownerid, args[2])
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定完成"))
			srv, _ := db.Find(grpid)
			msg := "群号：" + args[0] + "\n"
			msg += "腐竹qq：" + args[1] + "\n"
			msg += "服务器1数据：" + "\n" + "\t服务器名：" + srv.Servers[0].ServerName + "\n"
			msg += "\tgameid：" + srv.Servers[0].Gameid + "\n\tserverid：" + srv.Servers[0].Serverid
			Txt2Img(ctx, msg)
		})

	engine.OnPrefix(".绑服", zero.OnlyGroup, ServerOwnerPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			gids := strings.Split(ctx.State["args"].(string), " ")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在绑定..."))
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			var wg sync.WaitGroup
			for _, v := range gids {
				wg.Add(1)
				go func(s string) {
					err := db.AddServer(ctx.Event.GroupID, s)
					if err != nil {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：绑定 ", s, " 时发生错误", err))
					}
					wg.Done()
				}(v)
			}
			wg.Wait()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定结束"))
		})

	engine.OnPrefix(".changeowner", zero.OnlyGroup, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在修改服主信息..."))
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			grpid, _ := strconv.ParseInt(args[0], 10, 64)
			ownerid, _ := strconv.ParseInt(args[1], 10, 64)
			err = db.ChangeOwner(grpid, ownerid)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("修改成功"))
		})

	engine.OnPrefixGroup([]string{".addadmin", ".添加管理"}, zero.OnlyGroup, ServerOwnerPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			qids := strings.Split(ctx.State["args"].(string), " ")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在添加..."))
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			for _, v := range qids {
				qid, _ := strconv.ParseInt(v, 10, 64)
				err = db.AddAdmin(ctx.Event.GroupID, qid)
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：添加", v, "为管理时发生错误：", err))
				}
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("添加结束"))
		})

	engine.OnPrefixGroup([]string{".setalias", ".别名"}, zero.OnlyGroup, ServerOwnerPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			if len(args) != 2 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：非法参数！"))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("将", args[0], "的别名设置为", args[1]))
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			db.SetAlias(ctx.Event.GroupID, args[0], args[1])
		})

	// .kick name [reason]
	engine.OnPrefix(".kick", zero.OnlyGroup, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在踢出：", args[0]))
			if len(args) < 2 {
				args = append(args, "Admin Kick")
			}
			// 踢出理由转为繁体
			args[1] = S2tw(args[1])
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
			}
			data, err := db.Find(ctx.Event.GroupID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未绑定服务器"))
					return
				}
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			var wg sync.WaitGroup
			var reasons []string
			for _, v := range data.Servers {
				wg.Add(1)
				go func(v bf1model.Server) {
					var mu sync.Mutex
					srv := bf1rsp.NewServer(v.Serverid, v.Gameid, v.PGid)
					pid, err := bf1api.GetPersonalID(args[0])
					defer wg.Done()
					if err != nil {
						if v.NameInGroup != "" {
							mu.Lock()
							reasons = append(reasons, "在 "+v.NameInGroup+" 踢出失败："+err.Error())
							mu.Unlock()
							return
						}
						mu.Lock()
						reasons = append(reasons, "在 "+v.ServerName+" 踢出失败："+err.Error())
						mu.Unlock()
						return
					}
					reason, err := srv.Kick(pid, args[1])
					if err != nil {
						if v.NameInGroup != "" {
							mu.Lock()
							reasons = append(reasons, "在 "+v.NameInGroup+" 踢出失败："+err.Error())
							mu.Unlock()
							return
						}
						mu.Lock()
						reasons = append(reasons, "在 "+v.ServerName+" 踢出失败："+err.Error())
						mu.Unlock()
						return
					}
					if v.NameInGroup != "" {
						mu.Lock()
						reasons = append(reasons, "在 "+v.NameInGroup+" 踢出成功："+reason)
						mu.Unlock()
						return
					}
					mu.Lock()
					reasons = append(reasons, "在 "+v.NameInGroup+" 踢出成功："+reason)
					mu.Unlock()
				}(v)
			}
			wg.Wait()
			msg := "踢出 " + args[0] + "：\n"
			for _, v := range reasons {
				msg += "\t" + v + "\n"
			}
			Txt2Img(ctx, msg)
		})

	// .ban srv id
	engine.OnPrefix(".ban", zero.OnlyGroup, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			if len(args) < 2 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：非法参数！"))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在将 ", args[1], " 加入到 ", args[0], " 的ban列中"))
			db, cl, err := OpenServerDB()
			defer cl()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			s, err := db.GetServer(args[0], ctx.Event.GroupID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：找不到别名为 ", args[0], " 的服务器"))
					return
				}
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			srv := bf1rsp.NewServer(s.Serverid, s.Gameid, s.PGid)
			var wg sync.WaitGroup
			wg.Add(2)
			var pidchan chan string
			go func() {
				pid, err := bf1api.GetPersonalID(args[1])
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
					return
				}
				pidchan <- pid
				wg.Done()
			}()
			go func() {
				err := srv.Ban(<-pidchan)
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
					return
				}
			}()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已将 ", args[1], " 加入 ", args[0], " 的ban列"))
		})

	// .unban srv id
	engine.OnPrefix(".unban", zero.OnlyGroup, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			if len(args) < 2 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：非法参数！"))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在将 ", args[1], " 从 ", args[0], " 的ban列中删除"))
			db, cl, err := OpenServerDB()
			defer cl()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			s, err := db.GetServer(args[0], ctx.Event.GroupID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：找不到别名为 ", args[0], " 的服务器"))
					return
				}
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			srv := bf1rsp.NewServer(s.Serverid, s.Gameid, s.PGid)
			var wg sync.WaitGroup
			wg.Add(2)
			var pidchan chan string
			go func() {
				pid, err := bf1api.GetPersonalID(args[1])
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
					return
				}
				pidchan <- pid
				wg.Done()
			}()
			go func() {
				err := srv.Unban(<-pidchan)
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
					return
				}
			}()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已将 ", args[1], " 在 ", args[0], " 解封"))
		})

	engine.OnShell("自动踢出", autokick{}, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			info := ctx.State["flag"].(*autokick)
			// 检查参数
			if info.Server == "" {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：未指定服务器"))
				return
			}
			// 未指定等级
			if info.Rank == 0 {
				if info.Ping == 0 {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：等级和ping需要至少指定一个"))
					return
				}
				info.Rank = 151
			}
			// 未指定ping
			if info.Ping == 0 {
				info.Ping = 9999
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("将在 ", info.Server, " 踢出等级大于", info.Rank, "或ping高于", info.Ping, "的玩家"))
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
			var wg sync.WaitGroup
			wg.Add(1)
			srv := bf1rsp.NewServer(s.Serverid, s.Gameid, s.PGid)
			c := cron.New()
			go func() {
				c.AddFunc("@every 60s", func() {
					data, err := bf1api.ReturnJson(bf1api.EASBAPI, "POST", &post)
					if err != nil {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
						return
					}
					gjson.Get(data, "data."+s.Gameid+".players").ForEach(func(_, value gjson.Result) bool {
						if value.Get("rank").Int() > int64(info.Rank) {
							ctx.SendChain(message.Text("正在踢出", value.Get("display_name"), "：等级过高(", value.Get("rank"), ")"))
							srv.Kick(strconv.FormatInt(value.Get("pid").Int(), 10), "Rank limit "+strconv.Itoa(info.Rank))
						}
						if value.Get("ping").Int() > int64(info.Ping) {
							ctx.SendChain(message.Text("正在踢出", value.Get("display_name"), "：ping值过高(", value.Get("ping"), ")"))
							srv.Kick(strconv.FormatInt(value.Get("pid").Int(), 10), "Ping limit "+strconv.Itoa(info.Ping))
						}
						return true
					})
				})
				c.Start()
			}()
			next := zero.NewFutureEvent("message", 999, false, zero.OnlyGroup, ServerAdminPermission, zero.RegexRule(`^关闭自动踢出$`), zero.CheckGroup(ctx.Event.GroupID))
			recv, cancle := next.Repeat()
			defer cancle()
			for r := range recv {
				c.Stop()
				wg.Done()
				ctx.SendChain(message.Reply(r.Event.MessageID), message.Text("正在关闭自动踢出..."))
				ctx.SendChain(message.Text(info.Server, " 的自动踢出已关闭"))
				return
			}
			wg.Wait()
		})
}
