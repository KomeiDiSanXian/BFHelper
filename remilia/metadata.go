package remilia

import "github.com/KomeiDiSanXian/remilia/plugin"

func newMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Author:      "KomeiDiSanXian",
		Description: "战地系列游戏辅助插件 — 玩家战绩、服务器管理、联ban查询",
		Category:    "游戏工具",
		Tags:        []string{"战地", "Battlefield", "BF1", "BFV", "游戏", "服务器管理"},
		HelpText: `BFHelper — 战地系列游戏辅助插件

/login                   — 使用配置的 EA 凭据登录（仅超级管理员）
/logout                  — 登出 EA 账号（仅超级管理员）
/status                  — 查看连接状态

/bind <name>             — 绑定玩家名到 QQ
/unbind                  — 解绑
/stats [name]            — 查询玩家战绩（BTR）
/weapons [name]          — 查询武器数据
/vehicles [name]         — 查询载具数据
/recent [name]           — 查询最近战绩
/exchange                — 查询本期交换皮肤
/campaign                — 查询本期行动包
/cb [name]               — 查询联ban信息

/group create            — 创建服务器群组（群主）
/group delete            — 删除服务器群组（服主）
/group owner <qq>        — 更换服主（服主）
/group bind <gid...>     — 绑定服务器到群组（服主）
/group unbind <gid>      — 解绑服务器（服主）
/group admin add <qq...> — 添加群组管理员（服主）
/group admin rm <qq>     — 删除群组管理员（服主）

/admin kick <player>         — 踢出玩家（管理员）
/admin ban [a] <player>      — 封禁玩家（管理员）
/admin unban [a] <player>    — 解封玩家（管理员）
/admin banall <player>       — 全服封禁（管理员）
/admin unbanall <player>     — 全服解封（管理员）
/admin cm <a> [idx]          — 切换地图（管理员）
/admin qm <a>                — 查看地图池（管理员）
/admin move <p> <1|2> [a]    — 移动玩家到队伍（管理员）
/admin sadmin add|rm <p> [a] — 添加/移除服务器管理员（管理员）
/admin vip add|rm <p> [a]    — 添加/移除 VIP（管理员）
/admin logs <a> [n]          — 查询服务器操作日志（管理员）

/logs [起始] [结束] [类型]     — 查询本群操作日志，支持日期时间范围（管理员）

/blaze watch <gid>       — 监听服务器（管理员）
/blaze unwatch <gid>     — 取消监听（管理员）
/blaze list              — 查看监听列表（管理员）`,
	}
}
