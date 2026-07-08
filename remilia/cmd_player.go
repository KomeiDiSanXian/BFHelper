package remilia

import (
	stdctx "context"
	"fmt"
	"strconv"
	"strings"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/core"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/public"
	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/platform"
	"gorm.io/gorm"
)

func (p *Plugin) handleBind(ctx *eventctx.Context) error {
	parsed := ctx.GetParsedCommand()
	name := parsed.GetString("name")
	if name == "" {
		ctx.Reply(platform.TextMessage("用法: /bf bind <玩家名>"))
		return nil
	}

	qqID := p.getUserID(ctx)
	_, err := p.Storage.GetPlayerByQQ(qqID)
	if err == nil {
		existed, _ := p.Storage.GetPlayerByQQ(qqID)
		_ = p.Storage.DeletePlayer(qqID)
		ctx.Reply(platform.TextMessage(fmt.Sprintf("已将 %s 解绑，重新绑定为 %s", existed.PersonaName, name)))
	} else if err != gorm.ErrRecordNotFound {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("数据库错误: %v", err)))
		return nil
	}

	if err := p.Storage.CreatePlayer(&core.Player{QQID: qqID, PersonaName: name}); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("绑定失败: %v", err)))
		return nil
	}
	ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ 已绑定 %s", name)))
	return nil
}

func (p *Plugin) handleUnbind(ctx *eventctx.Context) error {
	qqID := p.getUserID(ctx)
	if err := p.Storage.DeletePlayer(qqID); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("解绑失败: %v", err)))
		return nil
	}
	ctx.Reply(platform.TextMessage("✅ 已解绑"))
	return nil
}

func (p *Plugin) resolvePersonaName(ctx *eventctx.Context, name string) (string, error) {
	if name != "" {
		return name, nil
	}
	qqID := p.getUserID(ctx)
	player, err := p.Storage.GetPlayerByQQ(qqID)
	if err != nil {
		return "", core.ErrPlayerNotFound
	}
	return player.PersonaName, nil
}

func (p *Plugin) handleStats(ctx *eventctx.Context) error {
	parsed := ctx.GetParsedCommand()
	name, err := p.resolvePersonaName(ctx, parsed.GetString("name"))
	if err != nil {
		ctx.Reply(platform.TextMessage(err.Error()))
		return nil
	}

	ctx.Reply(platform.TextMessage("少女折寿中..."))

	stat, err := p.BTRClient.GetStats(name)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询失败: %v", err)))
		return nil
	}

	msg := public.FormatStat(stat, name)
	ctx.Reply(platform.TextMessage(msg))
	return nil
}

func (p *Plugin) handleWeapons(ctx *eventctx.Context) error {
	parsed := ctx.GetParsedCommand()
	name, err := p.resolvePersonaName(ctx, parsed.GetString("name"))
	if err != nil {
		ctx.Reply(platform.TextMessage(err.Error()))
		return nil
	}

	if !p.Client.IsLoggedIn() {
		ctx.Reply(platform.TextMessage("需要先登录 EA 账号才能查询武器数据"))
		return nil
	}

	ctx.Reply(platform.TextMessage("少女折寿中..."))

	pid, err := p.resolvePersonaID(ctx.Context(), name)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("获取玩家信息失败: %v", err)))
		return nil
	}

	ws, err := public.GetWeapons(pid, public.SpartaAPIURL, p.Client.GatewaySession(), core.WeaponALL)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询武器失败: %v", err)))
		return nil
	}

	msg := public.FormatWeapons(ws, name)
	ctx.Reply(platform.TextMessage(msg))
	return nil
}

func (p *Plugin) handleVehicles(ctx *eventctx.Context) error {
	parsed := ctx.GetParsedCommand()
	name, err := p.resolvePersonaName(ctx, parsed.GetString("name"))
	if err != nil {
		ctx.Reply(platform.TextMessage(err.Error()))
		return nil
	}

	if !p.Client.IsLoggedIn() {
		ctx.Reply(platform.TextMessage("需要先登录 EA 账号才能查询载具数据"))
		return nil
	}

	ctx.Reply(platform.TextMessage("少女折寿中..."))

	pid, err := p.resolvePersonaID(ctx.Context(), name)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("获取玩家信息失败: %v", err)))
		return nil
	}

	vs, err := public.GetVehicles(pid, public.SpartaAPIURL, p.Client.GatewaySession())
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询载具失败: %v", err)))
		return nil
	}

	msg := public.FormatVehicles(vs, name)
	ctx.Reply(platform.TextMessage(msg))
	return nil
}

func (p *Plugin) handleRecent(ctx *eventctx.Context) error {
	parsed := ctx.GetParsedCommand()
	name, err := p.resolvePersonaName(ctx, parsed.GetString("name"))
	if err != nil {
		ctx.Reply(platform.TextMessage(err.Error()))
		return nil
	}

	ctx.Reply(platform.TextMessage("少女折寿中..."))

	entries, err := p.Bili22Client.GetRecent(name)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询失败: %v", err)))
		return nil
	}

	if len(entries) == 0 {
		ctx.Reply(platform.TextMessage("没有找到最近游玩记录"))
		return nil
	}

	msg := public.FormatRecent(entries, name)
	ctx.Reply(platform.TextMessage(msg))
	return nil
}

func (p *Plugin) handleExchange(ctx *eventctx.Context) error {
	if !p.Client.IsLoggedIn() {
		ctx.Reply(platform.TextMessage("需要先登录 EA 账号才能查询交换信息"))
		return nil
	}

	exchange, err := public.GetExchange(p.Client.GatewaySession())
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询失败: %v", err)))
		return nil
	}

	var msg strings.Builder
	for category, skins := range exchange {
		msg.WriteString(category + ":\n")
		for _, skin := range skins {
			msg.WriteString("  " + skin + "\n")
		}
	}
	ctx.Reply(platform.TextMessage(msg.String()))
	return nil
}

func (p *Plugin) handleCampaign(ctx *eventctx.Context) error {
	if !p.Client.IsLoggedIn() {
		ctx.Reply(platform.TextMessage("需要先登录 EA 账号才能查询行动包"))
		return nil
	}

	pack, err := public.GetCampaignPacks(p.Client.GatewaySession())
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询失败: %v", err)))
		return nil
	}

	msg := fmt.Sprintf("行动名：%s\n剩余时间：%.2f 天\n箱子重置：%.2f 小时\n地图：%s 与 %s\n简介：%s",
		pack.Name, float64(pack.RemainTime)/60/24, float64(pack.ResetTime)/60,
		pack.Op1Name, pack.Op2Name, pack.Desc)
	ctx.Reply(platform.TextMessage(msg))
	return nil
}

func (p *Plugin) handleCheaterCheck(ctx *eventctx.Context) error {
	parsed := ctx.GetParsedCommand()
	name, err := p.resolvePersonaName(ctx, parsed.GetString("name"))
	if err != nil {
		ctx.Reply(platform.TextMessage(err.Error()))
		return nil
	}

	ctx.Reply(platform.TextMessage("少女折寿中..."))

	info := public.CheckCheater(name)
	msg := fmt.Sprintf("ID: %s\nEAC: %s\n案件: %s\nBFBan: %s\n",
		name, info.EAC.Status, info.EAC.URL, info.BFBan.Status)
	if info.BFBan.IsCheater && info.BFBan.URL != "" {
		msg += "案件: " + info.BFBan.URL + "\n"
	}
	ctx.Reply(platform.TextMessage(msg))
	return nil
}

func (p *Plugin) resolvePersonaID(stdCtx stdctx.Context, name string) (string, error) {
	persona, err := p.Client.LookupPersonaByName(stdCtx, name)
	if err != nil {
		// Fallback: check local DB
		player, err2 := p.Storage.GetPlayerByName(name)
		if err2 != nil || player.PersonaID == "" {
			return "", err
		}
		return player.PersonaID, nil
	}
	return strconv.FormatInt(persona.PersonaID, 10), nil
}
