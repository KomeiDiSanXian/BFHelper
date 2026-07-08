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

	_, err := p.Storage.GetGroup(gid)
	if err == nil {
		ctx.Reply(platform.TextMessage("本群已存在服务器群组"))
		return nil
	}

	grp := &core.Group{GroupID: gid, OwnerQQ: qqID}
	if err := p.Storage.CreateGroup(grp); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("创建失败: %v", err)))
		return nil
	}
	ctx.Reply(platform.TextMessage("✅ 服务器群组已创建"))
	return nil
}

func (p *Plugin) handleGroupDelete(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)

	if err := p.Storage.DeleteGroup(gid); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("删除失败: %v", err)))
		return nil
	}
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
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf group owner <QQ号>"))
		return nil
	}

	newOwner, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		ctx.Reply(platform.TextMessage("无效的 QQ 号"))
		return nil
	}

	group, err := p.Storage.GetGroup(gid)
	if err != nil {
		ctx.Reply(platform.TextMessage("群组不存在"))
		return nil
	}

	group.OwnerQQ = newOwner
	if err := p.Storage.UpdateGroup(group); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("更新失败: %v", err)))
		return nil
	}
	ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ 服主已更换为 %d", newOwner)))
	return nil
}

func (p *Plugin) handleGroupBind(ctx *eventctx.Context) error {
	if !ctx.GetPlatformEvent().Chat().IsGroup {
		ctx.Reply(platform.TextMessage("该命令仅支持群聊"))
		return nil
	}
	gid := p.getGroupID(ctx)

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf group bind <gameId1> <gameId2> ..."))
		return nil
	}

	gameIDs := args[3:]
	if len(gameIDs) == 0 {
		ctx.Reply(platform.TextMessage("请指定至少一个 GameID"))
		return nil
	}

	success := 0
	for _, gidStr := range gameIDs {
		sv := &core.Server{GameID: gidStr, GroupID: gid}
		if err := p.Storage.AddServerToGroup(gid, sv); err != nil {
			continue
		}
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
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf group unbind <gameId>"))
		return nil
	}

	if err := p.Storage.RemoveServerFromGroup(gid, args[3]); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("解绑失败: %v", err)))
		return nil
	}
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
	if len(args) < 5 {
		ctx.Reply(platform.TextMessage("用法: /bf group admin add <QQ号1> <QQ号2> ..."))
		return nil
	}

	success := 0
	for _, qqStr := range args[4:] {
		qqID, err := strconv.ParseInt(qqStr, 10, 64)
		if err != nil {
			continue
		}
		if err := p.Storage.AddAdminToGroup(gid, qqID); err != nil {
			continue
		}
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
	if len(args) < 5 {
		ctx.Reply(platform.TextMessage("用法: /bf group admin rm <QQ号>"))
		return nil
	}

	qqID, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		ctx.Reply(platform.TextMessage("无效的 QQ 号"))
		return nil
	}

	if err := p.Storage.RemoveAdminFromGroup(gid, qqID); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("删除失败: %v", err)))
		return nil
	}
	ctx.Reply(platform.TextMessage("✅ 管理员已删除"))
	return nil
}
