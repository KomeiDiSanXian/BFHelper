package bfhelper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/rsp"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
)

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

	// .map srv mapid
	engine.OnPrefixGroup([]string{".map", ".maplist", ".切图"}, zero.OnlyGroup, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			if len(args) < 2 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：非法参数！"))
				return
			}
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
			maps, err := srv.GetMaps()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			var txt string
			for i, m := range *maps {
				txt += fmt.Sprintf("%d\t%s\t%s", i, m.ModeName, m.MapName)
			}
			txt += "--------\n请在30s内回复序号以进行切图"
			Txt2Img(ctx, txt)
			next := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^\d{1,2}$`), zero.OnlyGroup, ctx.CheckSession())
			recv, cancle := next.Repeat()
			defer cancle()
			timeout := time.NewTimer(30 * time.Second)
			for {
				select {
				case c := <-recv:
					i := c.Event.Message.String()
					idx, _ := strconv.Atoi(i)
					ctx.SendChain(message.Reply(c.Event.MessageID), message.Text("正在切换", args[0], "到", (*maps)[idx].ModeName, "\t", (*maps)[idx].MapName))
					err := srv.ChangeMap(idx)
					if err != nil {
						ctx.SendChain(message.Reply(c.Event.MessageID), message.Text("ERR：", err))
						return
					}
					ctx.SendChain(message.Reply(c.Event.MessageID), message.Text("切图完成"))
					return
				case <-timeout.C:
					return
				}
			}
		})
}
