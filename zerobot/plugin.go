// Package zerobot 提供 BFHelper 的 ZeroBot 插件适配器（公共功能）。
//
// 本适配器仅注册不依赖 BattlefieldAPI 私有库的公共命令：
//   - .bind / .绑定 — 账号绑定
//   - .战绩 / .stats — BTR 战绩查询
//   - .最近 / .recent — 最近游玩记录
//   - .cb — 联ban查询
//
// 需要 EA 登录的命令（武器、载具、交换、行动包）不被注册，
// 如需使用请通过 remilia 模式或自行补充 BattlefieldAPI 依赖。
//
// 导入方式：
//
//	import _ "github.com/KomeiDiSanXian/BFHelper/zerobot"
//
// 然后使用 ZeroBot 标准方式启动。
package zerobot

import (
	"fmt"
	"strings"
	"sync"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/core"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/public"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var store core.Storage

var (
	dbOnce sync.Once
	btr    = public.NewBTRClient()
	bili   = public.NewBili22Client()
)

func init() {
	dbOnce.Do(initDB)

	// ── 绑定 ──────────────────────────────────────────────────

	zero.OnPrefixGroup([]string{".绑定", ".bind"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["args"].(string)
		if id == "" {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败: 空 id"))
			return
		}
		player, err := store.GetPlayerByQQ(int64(ctx.Event.UserID))
		if err == nil {
			if player.PersonaName == id {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("笨蛋! 你现在绑的就是这个 id"))
				return
			}
			_ = store.DeletePlayer(int64(ctx.Event.UserID))
		}
		if err := store.CreatePlayer(&core.Player{QQID: int64(ctx.Event.UserID), PersonaName: id}); err != nil {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定失败: 数据库错误"))
			return
		}
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("绑定 ", id, " 成功"))
	})

	// ── 战绩（BTR）────────────────────────────────────────────

	zero.OnRegex(`^\. *1?战绩 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		name := resolveName(ctx)
		if name == "" {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("未绑定账号，请使用 .bind <id> 绑定"))
			return
		}
		ctx.Send("少女折寿中...")
		stat, err := btr.GetStats(name)
		if err != nil {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("查询失败: ", err.Error()))
			return
		}
		ctx.SendChain(message.Text(public.FormatStat(stat, name)))
	})

	// ── 最近战绩（Bili22）─────────────────────────────────────

	zero.OnRegex(`^\. *1?最近 *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		name := resolveName(ctx)
		if name == "" {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("未绑定账号，请使用 .bind <id> 绑定"))
			return
		}
		ctx.Send("少女折寿中...")
		entries, err := bili.GetRecent(name)
		if err != nil {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("查询失败: ", err.Error()))
			return
		}
		if len(entries) == 0 {
			ctx.SendChain(message.Text("没有找到最近游玩记录"))
			return
		}
		ctx.SendChain(message.Text(public.FormatRecent(entries, name)))
	})

	// ── 联ban查询 ─────────────────────────────────────────────

	zero.OnRegex(`^\. *1?cb *(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		name := resolveName(ctx)
		if name == "" {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("未绑定账号，请使用 .bind <id> 绑定"))
			return
		}
		ctx.Send("少女折寿中...")
		info := public.CheckCheater(name)
		msg := fmt.Sprintf("ID: %s\nEAC: %s\n案件: %s\nBFBan: %s\n",
			name, info.EAC.Status, info.EAC.URL, info.BFBan.Status)
		if info.BFBan.IsCheater && info.BFBan.URL != "" {
			msg += "案件: " + info.BFBan.URL + "\n"
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(msg))
	})
}

func initDB() {
	db, err := gorm.Open(sqlite.Open("data/battlefield.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("zerobot: failed to open database: %v", err))
	}
	if err := db.AutoMigrate(&core.Player{}); err != nil {
		panic(fmt.Sprintf("zerobot: failed to migrate database: %v", err))
	}
	store = core.NewStorage(db)
}

func resolveName(ctx *zero.Ctx) string {
	id := ctx.State["regex_matched"].([]string)[1]
	if id == "" {
		player, err := store.GetPlayerByQQ(int64(ctx.Event.UserID))
		if err != nil {
			return ""
		}
		return player.PersonaName
	}
	return strings.TrimSpace(id)
}
