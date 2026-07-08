package remilia

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/core"
	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/platform"
)

func (p *Plugin) requireLoggedIn(ctx *eventctx.Context) bool {
	if !p.Client.IsLoggedIn() {
		ctx.Reply(platform.TextMessage("需要先登录 EA 账号，请使用 /bf login"))
		return false
	}
	return true
}

func (p *Plugin) getGroupID(ctx *eventctx.Context) int64 {
	return parseUserID(ctx.GetPlatformEvent().Chat().ID)
}

func (p *Plugin) getUserID(ctx *eventctx.Context) int64 {
	return parseUserID(ctx.GetPlatformEvent().Sender().ID)
}

func (p *Plugin) requireGroupExists(ctx *eventctx.Context) (int64, bool) {
	gid := p.getGroupID(ctx)
	_, err := p.Storage.GetGroup(gid)
	if err != nil {
		ctx.Reply(platform.TextMessage("本群未创建服务器群组，请先使用 /bf group create"))
		return 0, false
	}
	return gid, true
}

func (p *Plugin) checkOwner(ctx *eventctx.Context, groupID int64) bool {
	qqID := p.getUserID(ctx)
	if p.Storage.IsServerOwner(groupID, qqID) {
		return true
	}
	ctx.Reply(platform.TextMessage("权限不足：需要服主权限"))
	return false
}

func (p *Plugin) checkAdmin(ctx *eventctx.Context, groupID int64) bool {
	qqID := p.getUserID(ctx)
	if p.Storage.IsServerAdmin(groupID, qqID) || p.Storage.IsServerOwner(groupID, qqID) {
		return true
	}
	ctx.Reply(platform.TextMessage("权限不足：需要服务器管理员权限"))
	return false
}

func (p *Plugin) handleAdminKick(ctx *eventctx.Context) error {
	if !p.requireLoggedIn(ctx) {
		return nil
	}
	gid, ok := p.requireGroupExists(ctx)
	if !ok {
		return nil
	}
	if !p.checkAdmin(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf admin kick <玩家名>"))
		return nil
	}
	playerName := args[3]

	group, _ := p.Storage.GetGroup(gid)
	results := make([]string, 0, len(group.Servers))
	for _, sv := range group.Servers {
		persona, err := p.Client.LookupPersonaByName(ctx.Context(), playerName)
		if err != nil {
			results = append(results, fmt.Sprintf("%s: 查询玩家失败", sv.ServerName))
			continue
		}
		pidStr := strconv.FormatInt(persona.PersonaID, 10)
		if _, err := p.Client.KickPlayer(sv.GameID, pidStr, "Kicked by admin"); err != nil {
			results = append(results, fmt.Sprintf("%s: %v", sv.ServerName, err))
			continue
		}
		results = append(results, fmt.Sprintf("%s: ✅ 已踢出", sv.ServerName))
	}

	ctx.Reply(platform.TextMessage(strings.Join(results, "\n")))
	return nil
}

func (p *Plugin) handleAdminBan(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.Client.BanPlayer(gameID, pid)
		return err
	})
}

func (p *Plugin) handleAdminUnban(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.Client.UnbanPlayer(gameID, pid)
		return err
	})
}

func (p *Plugin) banUnbanCmd(ctx *eventctx.Context, action func(string, string) error) error {
	if !p.requireLoggedIn(ctx) {
		return nil
	}
	gid, ok := p.requireGroupExists(ctx)
	if !ok {
		return nil
	}
	if !p.checkAdmin(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf admin ban|unban [别名] <玩家名>"))
		return nil
	}

	var alias, playerName string
	if len(args) == 4 {
		playerName = args[3]
	} else if len(args) >= 5 {
		alias = args[3]
		playerName = args[4]
	}

	persona, err := p.Client.LookupPersonaByName(ctx.Context(), playerName)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询玩家失败: %v", err)))
		return nil
	}
	pidStr := strconv.FormatInt(persona.PersonaID, 10)

	group, _ := p.Storage.GetGroup(gid)
	var servers []core.Server
	if alias != "" {
		sv, err := p.Storage.GetServerByAlias(gid, alias)
		if err != nil {
			ctx.Reply(platform.TextMessage(fmt.Sprintf("未找到别名为 %s 的服务器", alias)))
			return nil
		}
		servers = []core.Server{*sv}
	} else {
		servers = group.Servers
	}

	results := make([]string, 0, len(servers))
	for _, sv := range servers {
		if err := action(sv.GameID, pidStr); err != nil {
			results = append(results, fmt.Sprintf("%s: %v", sv.ServerName, err))
		} else {
			results = append(results, fmt.Sprintf("%s: ✅", sv.ServerName))
		}
	}

	ctx.Reply(platform.TextMessage(strings.Join(results, "\n")))
	return nil
}

func (p *Plugin) handleAdminBanAll(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.Client.BanPlayer(gameID, pid)
		return err
	})
}

func (p *Plugin) handleAdminUnbanAll(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.Client.UnbanPlayer(gameID, pid)
		return err
	})
}

func (p *Plugin) handleAdminChangeMap(ctx *eventctx.Context) error {
	if !p.requireLoggedIn(ctx) {
		return nil
	}
	gid, ok := p.requireGroupExists(ctx)
	if !ok {
		return nil
	}
	if !p.checkAdmin(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf admin cm <别名> [地图索引]"))
		return nil
	}

	alias := args[3]
	sv, err := p.Storage.GetServerByAlias(gid, alias)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("未找到别名为 %s 的服务器", alias)))
		return nil
	}

	if len(args) >= 5 && args[4] != "" {
		mapIndex, err := strconv.Atoi(args[4])
		if err != nil {
			ctx.Reply(platform.TextMessage("地图索引必须为数字"))
			return nil
		}
		if _, err := p.Client.SwitchMap(sv.PGID, mapIndex); err != nil {
			ctx.Reply(platform.TextMessage(fmt.Sprintf("切换地图失败: %v", err)))
			return nil
		}
		ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ %s 地图已切换到索引 %d", alias, mapIndex)))
	} else {
		result, err := p.Client.GetServerDetails(sv.GameID)
		if err != nil {
			ctx.Reply(platform.TextMessage(fmt.Sprintf("查询地图池失败: %v", err)))
			return nil
		}
		rotation := result.Get("result.rotation").Array()
		if len(rotation) == 0 {
			ctx.Reply(platform.TextMessage("该服务器的地图池为空"))
			return nil
		}
		var msg strings.Builder
		msg.WriteString(fmt.Sprintf("%s 的地图池:\n", alias))
		for i, m := range rotation {
			msg.WriteString(fmt.Sprintf("%d. %s (%s)\n", i, m.Get("mapPrettyName").Str, m.Get("modePrettyName").Str))
		}
		ctx.Reply(platform.TextMessage(msg.String()))
	}
	return nil
}

func (p *Plugin) handleAdminQueryMaps(ctx *eventctx.Context) error {
	gid, ok := p.requireGroupExists(ctx)
	if !ok {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf admin qm <别名>"))
		return nil
	}

	alias := args[3]
	sv, err := p.Storage.GetServerByAlias(gid, alias)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("未找到别名为 %s 的服务器", alias)))
		return nil
	}

	result, err := p.Client.GetServerDetails(sv.GameID)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询失败: %v", err)))
		return nil
	}

	rotation := result.Get("result.rotation").Array()
	if len(rotation) == 0 {
		ctx.Reply(platform.TextMessage("该服务器的地图池为空"))
		return nil
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("%s 的地图池:\n", alias))
	for i, m := range rotation {
		msg.WriteString(fmt.Sprintf("%d. %s (%s)\n", i, m.Get("mapPrettyName").Str, m.Get("modePrettyName").Str))
	}
	ctx.Reply(platform.TextMessage(msg.String()))
	return nil
}
