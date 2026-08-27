// Package remilia 提供 BFHelper 的 Remilia 框架插件适配器。
package remilia

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/core"
	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/infra/logger"
	"github.com/KomeiDiSanXian/remilia/platform"
)

func (p *Plugin) requireLoggedIn(ctx *eventctx.Context) bool {
	if !p.IsLoggedIn() {
		ctx.Reply(platform.TextMessage("需要先登录 EA 账号，请使用 /login"))
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
	_, err := p.GetGroup(gid)
	if err != nil {
		ctx.Reply(platform.TextMessage("本群未创建服务器群组，请先使用 /group create"))
		return 0, false
	}
	return gid, true
}

func (p *Plugin) checkOwner(ctx *eventctx.Context, groupID int64) bool {
	qqID := p.getUserID(ctx)
	if p.isSuperUser(qqID) || p.IsServerOwner(groupID, qqID) {
		return true
	}
	ctx.Reply(platform.TextMessage("权限不足：需要服主权限"))
	return false
}

func (p *Plugin) checkAdmin(ctx *eventctx.Context, groupID int64) bool {
	qqID := p.getUserID(ctx)
	if p.isSuperUser(qqID) || p.IsServerAdmin(groupID, qqID) || p.IsServerOwner(groupID, qqID) {
		return true
	}
	ctx.Reply(platform.TextMessage("权限不足：需要服务器管理员权限"))
	return false
}

// isSuperUser 判断 QQ 是否为配置中的机器人超级管理员。
func (p *Plugin) isSuperUser(qqID int64) bool {
	return slices.Contains(p.SuperUsers, qqID)
}

// requireSuperUser 仅允许机器人超级管理员执行（登录/登出等敏感操作）。
func (p *Plugin) requireSuperUser(ctx *eventctx.Context) bool {
	if p.isSuperUser(p.getUserID(ctx)) {
		return true
	}
	ctx.Reply(platform.TextMessage("权限不足：需要机器人超级管理员权限"))
	return false
}

// logOp 记录一条管理操作日志，失败不影响主流程。
func (p *Plugin) logOp(ctx *eventctx.Context, gid int64, action, target, serverName, detail string, success bool) {
	entry := &core.OperationLog{
		GroupID:    gid,
		OperatorQQ: p.getUserID(ctx),
		Action:     action,
		Target:     target,
		ServerName: serverName,
		Detail:     detail,
		Success:    success,
		CreatedAt:  time.Now(),
	}
	if err := p.RecordOperationLog(entry); err != nil {
		logger.WithError(err).Warn("[BFHelper] failed to record operation log")
	}
}

// selectServers 解析命令目标服务器：指定别名则单台，否则为群组内全部服务器。
func (p *Plugin) selectServers(ctx *eventctx.Context, gid int64, alias string) ([]core.Server, bool) {
	if alias != "" {
		sv, err := p.GetServerByAlias(gid, alias)
		if err != nil {
			ctx.Reply(platform.TextMessage(fmt.Sprintf("未找到别名为 %s 的服务器", alias)))
			return nil, false
		}
		return []core.Server{*sv}, true
	}
	group, err := p.GetGroup(gid)
	if err != nil {
		ctx.Reply(platform.TextMessage("群组不存在"))
		return nil, false
	}
	return group.Servers, true
}

// actionName 从命令内容中提取动作名（如 /admin ban → ban）。
func actionName(ctx *eventctx.Context) string {
	fields := strings.Fields(ctx.GetMessageContent())
	if len(fields) >= 2 {
		return strings.ToLower(fields[1])
	}
	return "admin"
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
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /admin kick <玩家名>"))
		return nil
	}
	playerName := args[2]

	group, _ := p.GetGroup(gid)
	results := make([]string, 0, len(group.Servers))
	for _, sv := range group.Servers {
		persona, err := p.LookupPersonaByName(ctx.Context(), playerName)
		if err != nil {
			p.logOp(ctx, gid, "kick", playerName, sv.ServerName, "查询玩家失败", false)
			results = append(results, fmt.Sprintf("%s: 查询玩家失败", sv.ServerName))
			continue
		}
		pidStr := strconv.FormatInt(persona.PersonaID, 10)
		if _, err := p.KickPlayer(sv.GameID, pidStr, "Kicked by admin"); err != nil {
			p.logOp(ctx, gid, "kick", playerName, sv.ServerName, err.Error(), false)
			results = append(results, fmt.Sprintf("%s: %v", sv.ServerName, err))
			continue
		}
		p.logOp(ctx, gid, "kick", playerName, sv.ServerName, "", true)
		results = append(results, fmt.Sprintf("%s: ✅ 已踢出", sv.ServerName))
	}

	ctx.Reply(platform.TextMessage(strings.Join(results, "\n")))
	return nil
}

func (p *Plugin) handleAdminBan(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.BanPlayer(gameID, pid)
		return err
	})
}

func (p *Plugin) handleAdminUnban(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.UnbanPlayer(gameID, pid)
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
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /admin ban|unban [别名] <玩家名>"))
		return nil
	}

	var alias, playerName string
	if len(args) == 3 {
		playerName = args[2]
	} else if len(args) >= 4 {
		alias = args[2]
		playerName = args[3]
	}

	persona, err := p.LookupPersonaByName(ctx.Context(), playerName)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询玩家失败: %v", err)))
		return nil
	}
	pidStr := strconv.FormatInt(persona.PersonaID, 10)

	group, _ := p.GetGroup(gid)
	var servers []core.Server
	if alias != "" {
		sv, err := p.GetServerByAlias(gid, alias)
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
			p.logOp(ctx, gid, actionName(ctx), playerName, sv.ServerName, err.Error(), false)
			results = append(results, fmt.Sprintf("%s: %v", sv.ServerName, err))
		} else {
			p.logOp(ctx, gid, actionName(ctx), playerName, sv.ServerName, "", true)
			results = append(results, fmt.Sprintf("%s: ✅", sv.ServerName))
		}
	}

	ctx.Reply(platform.TextMessage(strings.Join(results, "\n")))
	return nil
}

func (p *Plugin) handleAdminBanAll(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.BanPlayer(gameID, pid)
		return err
	})
}

func (p *Plugin) handleAdminUnbanAll(ctx *eventctx.Context) error {
	return p.banUnbanCmd(ctx, func(gameID, pid string) error {
		_, err := p.UnbanPlayer(gameID, pid)
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
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /admin cm <别名> [地图索引]"))
		return nil
	}

	alias := args[2]
	sv, err := p.GetServerByAlias(gid, alias)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("未找到别名为 %s 的服务器", alias)))
		return nil
	}

	if len(args) >= 4 && args[3] != "" {
		mapIndex, err := strconv.Atoi(args[3])
		if err != nil {
			ctx.Reply(platform.TextMessage("地图索引必须为数字"))
			return nil
		}
		if _, err := p.SwitchMap(sv.PGID, mapIndex); err != nil {
			p.logOp(ctx, gid, "cm", fmt.Sprintf("%s(索引%d)", sv.ServerName, mapIndex), alias, err.Error(), false)
			ctx.Reply(platform.TextMessage(fmt.Sprintf("切换地图失败: %v", err)))
			return nil
		}
		p.logOp(ctx, gid, "cm", fmt.Sprintf("索引%d", mapIndex), alias, "", true)
		ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ %s 地图已切换到索引 %d", alias, mapIndex)))
	} else {
		result, err := p.GetServerDetails(sv.GameID)
		if err != nil {
			ctx.Reply(platform.TextMessage(fmt.Sprintf("查询地图池失败: %v", err)))
			return nil
		}
		rotation := result.Get("rotation").Array()
		if len(rotation) == 0 {
			ctx.Reply(platform.TextMessage("该服务器的地图池为空"))
			return nil
		}
		var msg strings.Builder
		fmt.Fprintf(&msg, "%s 的地图池:\n", alias)
		for i, m := range rotation {
			fmt.Fprintf(&msg, "%d. %s (%s)\n", i, m.Get("mapPrettyName").Str, m.Get("modePrettyName").Str)
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
	if !p.checkAdmin(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /admin qm <别名>"))
		return nil
	}

	alias := args[2]
	sv, err := p.GetServerByAlias(gid, alias)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("未找到别名为 %s 的服务器", alias)))
		return nil
	}

	result, err := p.GetServerDetails(sv.GameID)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询失败: %v", err)))
		return nil
	}

	rotation := result.Get("rotation").Array()
	if len(rotation) == 0 {
		ctx.Reply(platform.TextMessage("该服务器的地图池为空"))
		return nil
	}

	var msg strings.Builder
	fmt.Fprintf(&msg, "%s 的地图池:\n", alias)
	for i, m := range rotation {
		fmt.Fprintf(&msg, "%d. %s (%s)\n", i, m.Get("mapPrettyName").Str, m.Get("modePrettyName").Str)
	}
	ctx.Reply(platform.TextMessage(msg.String()))
	return nil
}

// handleAdminMove 将玩家移动到指定队伍。
// 用法: /admin move <玩家名> <1|2> [别名]
func (p *Plugin) handleAdminMove(ctx *eventctx.Context) error {
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
		ctx.Reply(platform.TextMessage("用法: /admin move <玩家名> <1|2> [别名]"))
		return nil
	}
	playerName, teamStr := args[2], args[3]
	teamID, err := strconv.Atoi(teamStr)
	if err != nil || (teamID != 1 && teamID != 2) {
		ctx.Reply(platform.TextMessage("队伍编号必须为 1 或 2"))
		return nil
	}
	alias := ""
	if len(args) >= 5 {
		alias = args[4]
	}

	persona, err := p.LookupPersonaByName(ctx.Context(), playerName)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询玩家失败: %v", err)))
		return nil
	}
	pidStr := strconv.FormatInt(persona.PersonaID, 10)

	servers, ok := p.selectServers(ctx, gid, alias)
	if !ok {
		return nil
	}

	results := make([]string, 0, len(servers))
	for _, sv := range servers {
		if _, err := p.MovePlayer(sv.GameID, pidStr, teamID); err != nil {
			p.logOp(ctx, gid, "move", playerName, sv.ServerName, err.Error(), false)
			results = append(results, fmt.Sprintf("%s: %v", sv.ServerName, err))
			continue
		}
		p.logOp(ctx, gid, "move", playerName, sv.ServerName, fmt.Sprintf("队伍 %d", teamID), true)
		results = append(results, fmt.Sprintf("%s: ✅ 已移动至队伍 %d", sv.ServerName, teamID))
	}
	ctx.Reply(platform.TextMessage(strings.Join(results, "\n")))
	return nil
}

// handleAdminServerAdmin 管理服务器 RSP 管理员。
// 用法: /admin sadmin add|rm <玩家名> [别名]
func (p *Plugin) handleAdminServerAdmin(ctx *eventctx.Context) error {
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
		ctx.Reply(platform.TextMessage("用法: /admin sadmin add|rm <玩家名> [别名]"))
		return nil
	}
	sub, playerName := args[2], args[3]
	if sub != "add" && sub != "rm" {
		ctx.Reply(platform.TextMessage("用法: /admin sadmin add|rm <玩家名> [别名]"))
		return nil
	}
	alias := ""
	if len(args) >= 5 {
		alias = args[4]
	}

	persona, err := p.LookupPersonaByName(ctx.Context(), playerName)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询玩家失败: %v", err)))
		return nil
	}
	pidStr := strconv.FormatInt(persona.PersonaID, 10)

	action := "sadmin_add"
	actionText := "添加服务器管理员"
	if sub == "rm" {
		action = "sadmin_rm"
		actionText = "移除服务器管理员"
	}

	servers, ok := p.selectServers(ctx, gid, alias)
	if !ok {
		return nil
	}

	results := make([]string, 0, len(servers))
	for _, sv := range servers {
		var err error
		if sub == "add" {
			_, err = p.AddServerAdmin(sv.GameID, pidStr)
		} else {
			_, err = p.RemoveServerAdmin(sv.GameID, pidStr)
		}
		if err != nil {
			p.logOp(ctx, gid, action, playerName, sv.ServerName, err.Error(), false)
			results = append(results, fmt.Sprintf("%s: %v", sv.ServerName, err))
			continue
		}
		p.logOp(ctx, gid, action, playerName, sv.ServerName, "", true)
		results = append(results, fmt.Sprintf("%s: ✅ %s", sv.ServerName, actionText))
	}
	ctx.Reply(platform.TextMessage(strings.Join(results, "\n")))
	return nil
}

// handleAdminVIP 管理服务器 VIP。
// 用法: /admin vip add|rm <玩家名> [别名]
func (p *Plugin) handleAdminVIP(ctx *eventctx.Context) error {
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
		ctx.Reply(platform.TextMessage("用法: /admin vip add|rm <玩家名> [别名]"))
		return nil
	}
	sub, playerName := args[2], args[3]
	if sub != "add" && sub != "rm" {
		ctx.Reply(platform.TextMessage("用法: /admin vip add|rm <玩家名> [别名]"))
		return nil
	}
	alias := ""
	if len(args) >= 5 {
		alias = args[4]
	}

	persona, err := p.LookupPersonaByName(ctx.Context(), playerName)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询玩家失败: %v", err)))
		return nil
	}
	pidStr := strconv.FormatInt(persona.PersonaID, 10)

	action := "vip_add"
	actionText := "添加 VIP"
	if sub == "rm" {
		action = "vip_rm"
		actionText = "移除 VIP"
	}

	servers, ok := p.selectServers(ctx, gid, alias)
	if !ok {
		return nil
	}

	results := make([]string, 0, len(servers))
	for _, sv := range servers {
		var err error
		if sub == "add" {
			_, err = p.AddServerVIP(sv.GameID, pidStr)
		} else {
			_, err = p.RemoveServerVIP(sv.GameID, pidStr)
		}
		if err != nil {
			p.logOp(ctx, gid, action, playerName, sv.ServerName, err.Error(), false)
			results = append(results, fmt.Sprintf("%s: %v", sv.ServerName, err))
			continue
		}
		p.logOp(ctx, gid, action, playerName, sv.ServerName, "", true)
		results = append(results, fmt.Sprintf("%s: ✅ %s", sv.ServerName, actionText))
	}
	ctx.Reply(platform.TextMessage(strings.Join(results, "\n")))
	return nil
}
