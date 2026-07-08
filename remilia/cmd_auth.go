package remilia

import (
	"fmt"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/core"
	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/platform"
)

func (p *Plugin) handleLogin(ctx *eventctx.Context) error {
	if p.Client.IsLoggedIn() {
		ctx.Reply(platform.TextMessage("已经登录，如需重新登录请先 /bf logout"))
		return nil
	}

	sid := p.Config.Account.SID
	remid := p.Config.Account.Remid
	if sid == "" || remid == "" {
		ctx.Reply(platform.TextMessage("未配置 EA 凭据，请在插件配置中设置 account.sid 和 account.remid"))
		return nil
	}

	eaGame := toEAGame(p.Config.Account.Game)
	if err := p.Client.Login(ctx.Context(), sid, remid, eaGame); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("登录失败: %v", err)))
		return nil
	}

	game, _ := core.ParseGame(p.Config.Account.Game)
	p.persistCredentials(sid, remid, game)

	ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ EA 登录成功（游戏: %s）", p.Config.Account.Game)))
	return nil
}

func (p *Plugin) handleLogout(ctx *eventctx.Context) error {
	p.Client.Logout()
	if p.Config.Account.SID == "" {
		_ = p.kv.Delete([]byte("ea_credentials"))
	}
	ctx.Reply(platform.TextMessage("✅ 已登出 EA 账号"))
	return nil
}

func (p *Plugin) handleStatus(ctx *eventctx.Context) error {
	msg := "BFHelper 状态\n"
	msg += "─────────────\n"

	if p.Client.IsLoggedIn() {
		msg += "EA 登录: ✅\n"
	} else {
		msg += "EA 登录: ❌\n"
	}

	db := p.Storage.DB()
	var playerCount, groupCount, watchCount int64
	db.Model(&core.Player{}).Count(&playerCount)
	db.Model(&core.Group{}).Count(&groupCount)
	db.Model(&core.BlazeWatch{}).Where("enabled = ?", true).Count(&watchCount)
	msg += fmt.Sprintf("绑定玩家: %d\n", playerCount)
	msg += fmt.Sprintf("服务器群组: %d\n", groupCount)
	msg += fmt.Sprintf("Blaze 监听: %d\n", watchCount)

	ctx.Reply(platform.TextMessage(msg))
	return nil
}

func parseUserID(idStr string) int64 {
	var id int64
	fmt.Sscanf(idStr, "%d", &id)
	return id
}
