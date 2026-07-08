package remilia

import (
	"fmt"
	"strings"

	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/platform"
)

func (p *Plugin) handleBlazeWatch(ctx *eventctx.Context) error {
	gid := p.getGroupID(ctx)

	_, err := p.Storage.GetGroup(gid)
	if err != nil {
		ctx.Reply(platform.TextMessage("请先创建服务器群组: /bf group create"))
		return nil
	}

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf blaze watch <gameId>"))
		return nil
	}

	gameID := args[3]
	if err := p.Storage.SetBlazeWatch(gid, gameID, true); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("设置失败: %v", err)))
		return nil
	}

	ctx.Reply(platform.TextMessage(fmt.Sprintf("✅ 已开始监听服务器 %s\n当该服务器有玩家加入/离开时将推送通知", gameID)))
	return nil
}

func (p *Plugin) handleBlazeUnwatch(ctx *eventctx.Context) error {
	gid := p.getGroupID(ctx)

	args := strings.Fields(ctx.GetMessageContent())
	if len(args) < 4 {
		ctx.Reply(platform.TextMessage("用法: /bf blaze unwatch <gameId>"))
		return nil
	}

	if err := p.Storage.RemoveBlazeWatch(gid, args[3]); err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("取消监听失败: %v", err)))
		return nil
	}

	ctx.Reply(platform.TextMessage("✅ 已取消监听"))
	return nil
}

func (p *Plugin) handleBlazeList(ctx *eventctx.Context) error {
	gid := p.getGroupID(ctx)

	watches, err := p.Storage.GetBlazeWatches(gid)
	if err != nil {
		ctx.Reply(platform.TextMessage(fmt.Sprintf("查询失败: %v", err)))
		return nil
	}

	if len(watches) == 0 {
		ctx.Reply(platform.TextMessage("本群没有正在监听的服务器\n使用 /bf blaze watch <gameId> 添加"))
		return nil
	}

	var msg strings.Builder
	msg.WriteString("📡 正在监听的服务器:\n")
	for _, w := range watches {
		msg.WriteString(fmt.Sprintf("  - %s\n", w.GameID))
	}
	ctx.Reply(platform.TextMessage(msg.String()))
	return nil
}
