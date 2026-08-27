package remilia

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/core"
	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/platform"
)

func (p *Plugin) handleGroupCreate(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)
	qqID := p.getUserID(ctx)

	_, err := p.GetGroup(gid)
	if err == nil {
		ctx.Reply(platform.TextMessage("本群已存在服务器群组"))
		return nil
	}

	grp := &core.Group{GroupID: gid, OwnerQQ: qqID}
	if err := p.CreateGroup(grp); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("创建失败: %v", err)))
		return nil
	}
	p.logOp(ctx, gid, "group_create", strconv.FormatInt(qqID, 10), "", "", true)
	ctx.Reply(platform.TextMessage("✅ 服务器群组已创建"))
	return nil
}

func (p *Plugin) handleGroupDelete(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)
	if _, err := p.GetGroup(gid); err != nil {
		ctx.Reply(platform.TextMessage("本群未创建服务器群组"))
		return nil
	}
	if !p.checkOwner(ctx, gid) {
		return nil
	}

	if err := p.DeleteGroup(gid); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("删除失败: %v", err)))
		return nil
	}
	p.logOp(ctx, gid, "group_delete", "", "", "", true)
	ctx.Reply(platform.TextMessage("✅ 服务器群组已删除"))
	return nil
}

func (p *Plugin) handleGroupOwner(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)
	if !p.checkOwner(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /group owner <QQ号>"))
		return nil
	}

	newOwner, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		ctx.Reply(platform.TextMessage("无效的 QQ 号"))
		return nil
	}

	group, err := p.GetGroup(gid)
	if err != nil {
		ctx.Reply(platform.TextMessage("群组不存在"))
		return nil
	}

	group.OwnerQQ = newOwner
	if err := p.UpdateGroup(group); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("更新失败: %v", err)))
		return nil
	}
	p.logOp(ctx, gid, "group_owner", strconv.FormatInt(newOwner, 10), "", "", true)
	ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ 服主已更换为 %d", newOwner)))
	return nil
}

func (p *Plugin) handleGroupBind(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	if !p.requireLoggedIn(ctx) {
		return nil
	}
	gid := p.getGroupID(ctx)
	if !p.checkOwner(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /group bind <gameId1> <gameId2> ..."))
		return nil
	}

	gameIDs := args[2:]
	if len(gameIDs) == 0 {
		ctx.Reply(platform.TextMessage("请指定至少一个 GameID"))
		return nil
	}

	success := 0
	for _, gidStr := range gameIDs {
		sv := &core.Server{GameID: gidStr, GroupID: gid}
		if info, err := p.GetServerInfo(ctx.Context(), gidStr); err == nil {
			sv.ServerID = info.ServerID
			sv.PGID = info.PGID
			sv.ServerName = info.ServerName
		}
		if err := p.AddServerToGroup(gid, sv); err != nil {
			continue
		}
		p.logOp(ctx, gid, "group_bind", gidStr, sv.ServerName, "", true)
		success++
	}

	ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ 成功绑定 %d/%d 个服务器", success, len(gameIDs))))
	return nil
}

func (p *Plugin) handleGroupUnbind(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)
	if !p.checkOwner(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 3 {
		ctx.Reply(platform.TextMessage("用法: /group unbind <gameId>"))
		return nil
	}

	if err := p.RemoveServerFromGroup(gid, args[2]); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("解绑失败: %v", err)))
		return nil
	}
	p.logOp(ctx, gid, "group_unbind", args[2], "", "", true)
	ctx.Reply(platform.TextMessage("✅ 服务器已解绑"))
	return nil
}

func (p *Plugin) handleGroupAdminAdd(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)
	if !p.checkOwner(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /group admin add <QQ号1> <QQ号2> ..."))
		return nil
	}

	success := 0
	for _, qqStr := range args[3:] {
		qqID, err := strconv.ParseInt(qqStr, 10, 64)
		if err != nil {
			continue
		}
		if err := p.AddAdminToGroup(gid, qqID); err != nil {
			continue
		}
		p.logOp(ctx, gid, "group_admin_add", strconv.FormatInt(qqID, 10), "", "", true)
		success++
	}

	ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ 成功添加 %d 个管理员", success)))
	return nil
}

func (p *Plugin) handleGroupAdminRemove(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)
	if !p.checkOwner(ctx, gid) {
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /group admin rm <QQ号>"))
		return nil
	}

	qqID, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		ctx.Reply(platform.TextMessage("无效的 QQ 号"))
		return nil
	}

	if err := p.RemoveAdminFromGroup(gid, qqID); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("删除失败: %v", err)))
		return nil
	}
	p.logOp(ctx, gid, "group_admin_rm", strconv.FormatInt(qqID, 10), "", "", true)
	ctx.Reply(platform.TextMessage("✅ 管理员已删除"))
	return nil
}
