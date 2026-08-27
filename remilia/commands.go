package remilia

import (
	"github.com/KomeiDiSanXian/remilia/command"
	eventctx "github.com/KomeiDiSanXian/remilia/core/context"
	"github.com/KomeiDiSanXian/remilia/plugin"
)

func (p *Plugin) registerCommands(ctx *plugin.SetupContext) {
	p.registerAuthCommands(ctx)
	p.registerPlayerCommands(ctx)
	p.registerGroupCommands(ctx)
	p.registerAdminCommands(ctx)
	p.registerBlazeCommands(ctx)
	p.registerLogCommands(ctx)
}

func (p *Plugin) onCommand(ctx *plugin.SetupContext, trigger string, def *command.Definition, handler eventctx.Handler) {
	ctx.OnCommandDefWith("", trigger, def, handler)
}

// ─── Auth ────────────────────────────────────────────────────────────────────

func (p *Plugin) registerAuthCommands(ctx *plugin.SetupContext) {
	p.onCommand(ctx, "/login",
		command.NewDef("login").Description("使用配置的 EA 凭据登录").Build(),
		p.handleLogin)

	p.onCommand(ctx, "/logout",
		command.NewDef("logout").Description("登出 EA 账号").Build(),
		p.handleLogout)

	p.onCommand(ctx, "/status",
		command.NewDef("status").Description("查看连接和登录状态").Build(),
		p.handleStatus)
}

// ─── Player ──────────────────────────────────────────────────────────────────

func (p *Plugin) registerPlayerCommands(ctx *plugin.SetupContext) {
	p.onCommand(ctx, "/bind",
		command.NewDef("bind").Description("绑定玩家名到 QQ").
			Arg("name", "玩家名称", true).Build(),
		p.handleBind)

	p.onCommand(ctx, "/unbind",
		command.NewDef("unbind").Description("解绑玩家名").Build(),
		p.handleUnbind)

	p.onCommand(ctx, "/stats",
		command.NewDef("stats").Description("查询玩家战绩（BTR）").
			Arg("name", "玩家名", false).Build(),
		p.handleStats)

	p.onCommand(ctx, "/weapons",
		command.NewDef("weapons").Description("查询武器数据").
			Arg("name", "玩家名", false).Build(),
		p.handleWeapons)

	p.onCommand(ctx, "/vehicles",
		command.NewDef("vehicles").Description("查询载具数据").
			Arg("name", "玩家名", false).Build(),
		p.handleVehicles)

	p.onCommand(ctx, "/recent",
		command.NewDef("recent").Description("查询最近战绩").
			Arg("name", "玩家名", false).Build(),
		p.handleRecent)

	p.onCommand(ctx, "/exchange",
		command.NewDef("exchange").Description("查询本期交换皮肤").Build(),
		p.handleExchange)

	p.onCommand(ctx, "/campaign",
		command.NewDef("campaign").Description("查询本期行动包").Build(),
		p.handleCampaign)

	p.onCommand(ctx, "/cb",
		command.NewDef("cb").Description("查询联ban信息").
			Arg("name", "玩家名", false).Build(),
		p.handleCheaterCheck)
}

// ─── Group ───────────────────────────────────────────────────────────────────

func (p *Plugin) registerGroupCommands(ctx *plugin.SetupContext) {
	p.onCommand(ctx, "/group create",
		command.NewDef("group").Description("管理服务器群组").
			SubCommand(command.NewDef("create").Description("创建服务器群组").Build()).
			Build(),
		p.handleGroupCreate)

	p.onCommand(ctx, "/group delete",
		command.NewDef("group").SubCommand(command.NewDef("delete").Description("删除服务器群组").Build()).Build(),
		p.handleGroupDelete)

	p.onCommand(ctx, "/group owner",
		command.NewDef("group").SubCommand(command.NewDef("owner").Description("更换服主").Build()).Build(),
		p.handleGroupOwner)

	p.onCommand(ctx, "/group bind",
		command.NewDef("group").SubCommand(command.NewDef("bind").Description("绑定服务器到群组").Build()).Build(),
		p.handleGroupBind)

	p.onCommand(ctx, "/group unbind",
		command.NewDef("group").SubCommand(command.NewDef("unbind").Description("解绑服务器").Build()).Build(),
		p.handleGroupUnbind)

	p.onCommand(ctx, "/group admin add",
		command.NewDef("group").SubCommand(command.NewDef("admin").Description("管理服务器管理员").
			SubCommand(command.NewDef("add").Description("添加服务器管理员").Build()).Build()).Build(),
		p.handleGroupAdminAdd)

	p.onCommand(ctx, "/group admin rm",
		command.NewDef("group").SubCommand(command.NewDef("admin").
			SubCommand(command.NewDef("rm").Description("删除服务器管理员").Build()).Build()).Build(),
		p.handleGroupAdminRemove)
}

// ─── Admin ───────────────────────────────────────────────────────────────────

func (p *Plugin) registerAdminCommands(ctx *plugin.SetupContext) {
	p.onCommand(ctx, "/admin kick",
		command.NewDef("admin").Description("服务器管理命令").
			SubCommand(command.NewDef("kick").Description("踢出玩家").Build()).Build(),
		p.handleAdminKick)

	p.onCommand(ctx, "/admin ban",
		command.NewDef("admin").SubCommand(command.NewDef("ban").Description("封禁玩家").Build()).Build(),
		p.handleAdminBan)

	p.onCommand(ctx, "/admin unban",
		command.NewDef("admin").SubCommand(command.NewDef("unban").Description("解封玩家").Build()).Build(),
		p.handleAdminUnban)

	p.onCommand(ctx, "/admin banall",
		command.NewDef("admin").SubCommand(command.NewDef("banall").Description("全服封禁").Build()).Build(),
		p.handleAdminBanAll)

	p.onCommand(ctx, "/admin unbanall",
		command.NewDef("admin").SubCommand(command.NewDef("unbanall").Description("全服解封").Build()).Build(),
		p.handleAdminUnbanAll)

	p.onCommand(ctx, "/admin cm",
		command.NewDef("admin").SubCommand(command.NewDef("cm").Description("切换地图").Build()).Build(),
		p.handleAdminChangeMap)

	p.onCommand(ctx, "/admin qm",
		command.NewDef("admin").SubCommand(command.NewDef("qm").Description("查看地图池").Build()).Build(),
		p.handleAdminQueryMaps)

	p.onCommand(ctx, "/admin move",
		command.NewDef("admin").SubCommand(command.NewDef("move").Description("将玩家移动到指定队伍").Build()).Build(),
		p.handleAdminMove)

	p.onCommand(ctx, "/admin sadmin add",
		command.NewDef("admin").SubCommand(command.NewDef("sadmin").Description("管理服务器管理员").
			SubCommand(command.NewDef("add").Description("添加服务器管理员").Build()).Build()).Build(),
		p.handleAdminServerAdmin)

	p.onCommand(ctx, "/admin sadmin rm",
		command.NewDef("admin").SubCommand(command.NewDef("sadmin").
			SubCommand(command.NewDef("rm").Description("移除服务器管理员").Build()).Build()).Build(),
		p.handleAdminServerAdmin)

	p.onCommand(ctx, "/admin vip add",
		command.NewDef("admin").SubCommand(command.NewDef("vip").Description("管理服务器 VIP").
			SubCommand(command.NewDef("add").Description("添加 VIP").Build()).Build()).Build(),
		p.handleAdminVIP)

	p.onCommand(ctx, "/admin vip rm",
		command.NewDef("admin").SubCommand(command.NewDef("vip").
			SubCommand(command.NewDef("rm").Description("移除 VIP").Build()).Build()).Build(),
		p.handleAdminVIP)

	p.onCommand(ctx, "/admin logs",
		command.NewDef("admin").SubCommand(command.NewDef("logs").Description("查询服务器操作日志").Build()).Build(),
		p.handleAdminLogs)
}

// registerLogCommands 注册操作日志查询命令。
func (p *Plugin) registerLogCommands(ctx *plugin.SetupContext) {
	p.onCommand(ctx, "/logs",
		command.NewDef("logs").Description("查询本群操作日志（支持日期时间范围）").Build(),
		p.handleLogs)
}

// ─── Blaze ───────────────────────────────────────────────────────────────────

func (p *Plugin) registerBlazeCommands(ctx *plugin.SetupContext) {
	p.onCommand(ctx, "/blaze watch",
		command.NewDef("blaze").Description("Blaze 实时推送管理").
			SubCommand(command.NewDef("watch").Description("开始监听指定服务器").Build()).Build(),
		p.handleBlazeWatch)

	p.onCommand(ctx, "/blaze unwatch",
		command.NewDef("blaze").SubCommand(command.NewDef("unwatch").Description("取消监听服务器").Build()).Build(),
		p.handleBlazeUnwatch)

	p.onCommand(ctx, "/blaze list",
		command.NewDef("blaze").SubCommand(command.NewDef("list").Description("查看监听列表").Build()).Build(),
		p.handleBlazeList)
}
