package remilia

import "github.com/KomeiDiSanXian/remilia/plugin"

func newMetadata() *plugin.Metadata {
	return &plugin.Metadata{
		Author:      "KomeiDiSanXian",
		Description: "战地系列游戏辅助插件 — 玩家战绩、服务器管理、联ban查询",
		Category:    "游戏工具",
		Tags:        []string{"战地", "Battlefield", "BF1", "BFV", "游戏", "服务器管理"},
		HelpText: `BFHelper — 战地系列游戏辅助插件

/bf login                   — 使用配置的 EA 凭据登录
/bf logout                  — 登出 EA 账号
/bf status                  — 查看连接状态

/bf bind <name>             — 绑定玩家名到 QQ
/bf unbind                  — 解绑
/bf stats [name]            — 查询玩家战绩（BTR）
/bf weapons [name]          — 查询武器数据
/bf vehicles [name]         — 查询载具数据
/bf recent [name]           — 查询最近战绩
/bf exchange                — 查询本期交换皮肤
/bf campaign                — 查询本期行动包
/bf cb [name]               — 查询联ban信息

/bf group create            — 创建服务器群组
/bf group delete            — 删除服务器群组
/bf group owner <qq>        — 更换服主
/bf group bind <gid...>     — 绑定服务器到群组
/bf group unbind <gid>      — 解绑服务器
/bf group admin add <qq...> — 添加服务器管理员
/bf group admin rm <qq>     — 删除服务器管理员

/bf admin kick <player>     — 踢出玩家
/bf admin ban [a] <player>  — 封禁玩家（指定别名）
/bf admin unban [a] <p>     — 解封玩家
/bf admin banall <player>   — 全服封禁
/bf admin unbanall <player> — 全服解封
/bf admin cm <a> [idx]      — 切换地图
/bf admin qm <a>            — 查看地图池

/bf blaze watch <gid>       — 监听服务器
/bf blaze unwatch <gid>     — 取消监听
/bf blaze list              — 查看监听列表`,
	}
}
