// Package bfhelper 战地玩家查询
package bfhelper

import (
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/handler"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
)

func init() {
	// QQ绑定ID
	global.Engine.OnPrefixGroup([]string{".绑定", ".bind"}).SetBlock(true).Handle(handler.BindAccountHandler)
	// bf1个人战绩
	global.Engine.OnRegex(`\. *1?战绩 *(.*)$`).SetBlock(true).Handle(handler.PlayerStatsHandler)
	// 武器查询，只展示前五个
	global.Engine.OnRegex(`^\. *1?武器 *(.*)$`).SetBlock(true).Handle(handler.PlayerWeaponHandler)
	// 最近战绩
	global.Engine.OnRegex(`^\. *1?最近 *(.*)$`).SetBlock(true).Handle(handler.PlayerRecentHandler)
	// 获取所有种类的载具信息
	global.Engine.OnRegex(`^\. *1?载具 *(.*)$`).SetBlock(true).Handle(handler.PlayerVehicleHandler)
	// 交换查询
	global.Engine.OnFullMatchGroup([]string{".交换", ".exchange"}).SetBlock(true).Handle(handler.BF1ExchangeHandler)
	// 行动包查询
	global.Engine.OnFullMatchGroup([]string{".行动", ".行动包", ".pack"}).SetBlock(true).Handle(handler.BF1OpreationPackHandler)
	// 查询玩家是否有案件
	global.Engine.OnRegex(`^\. *1?cb *(.*)$`).SetBlock(true).Handle(handler.PlayerBanInfoHandler)
}
