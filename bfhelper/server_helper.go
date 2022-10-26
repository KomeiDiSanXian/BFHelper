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
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			grpid, _ := strconv.ParseInt(args[0], 10, 64)
			ownerid, _ := strconv.ParseInt(args[1], 10, 64)
			err = db.Create(grpid, ownerid, args[3])
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

	engine.OnPrefix(".绑定服务器", zero.OnlyGroup, ServerOwnerPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			gids := strings.Split(ctx.State["args"].(string), " ")
			db, close, err := OpenServerDB()
			defer close()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
			}
			var wg sync.WaitGroup
			for _, v := range gids {
				wg.Add(1)
				go func(s string) {
					err = db.AddServer(ctx.Event.GroupID, s)
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：绑定 ", s, " 时发生错误", err))
					wg.Done()
				}(v)
			}
			wg.Wait()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定结束"))
		})

	engine.OnPrefixGroup([]string{".addadmin", ".添加管理"}, zero.OnlyGroup, ServerOwnerPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			qids := strings.Split(ctx.State["args"].(string), " ")
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
	// .kick name [reason]
	engine.OnPrefix(".kick", zero.OnlyGroup, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			args := strings.Split(ctx.State["args"].(string), " ")
			if len(args) < 2 {
				args = append(args, "Admin Kick")
			}
			// 踢出理由转为繁体
			args[1] = S2tw(args[1])
			// 理由过长
			if len(args[1]) > 21 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：理由过长！"))
				return
			}
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
					srv := bf1rsp.NewServer(v.Serverid, v.Gameid, v.PGid)
					pid, err := bf1api.GetPersonalID(args[0])
					defer wg.Done()
					if err != nil {
						if v.NameInGroup != "" {
							reasons = append(reasons, "在 "+v.NameInGroup+" 踢出失败："+err.Error())
							return
						}
						reasons = append(reasons, "在 "+v.ServerName+" 踢出失败："+err.Error())
						return
					}
					reason, err := srv.Kick(pid, args[1])
					if err != nil {
						if v.NameInGroup != "" {
							reasons = append(reasons, "在 "+v.NameInGroup+" 踢出失败："+err.Error())
							return
						}
						reasons = append(reasons, "在 "+v.ServerName+" 踢出失败："+err.Error())
						return
					}
					if v.NameInGroup != "" {
						reasons = append(reasons, "在 "+v.NameInGroup+" 踢出成功："+reason)
						return
					}
					reasons = append(reasons, "在 "+v.NameInGroup+" 踢出成功："+reason)
				}(v)
			}
			wg.Wait()
			msg := "踢出 " + args[0] + "：\n"
			for _, v := range reasons {
				msg += "\t" + v + "\n"
			}
			Txt2Img(ctx, msg)
		})
/*
	engine.OnShell("自动踢出", autokick{Rank: 150, Ping: 999}, ServerAdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			info := ctx.State["flag"].(*autokick)
			db, cl, err := OpenServerDB()
			defer cl()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			srv, err := db.Find(ctx.Event.GroupID)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR：", err))
				return
			}
			var post map[string][]string = make(map[string][]string)
			post["gameIds"] = []string{}
			bf1api.ReturnJson(bf1api.EASBAPI, "POST", &post)
		})*/
}
