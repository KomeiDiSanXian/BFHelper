// Package core 提供 BFHelper 核心库：共享类型、错误与存储层。
package core

import "fmt"

// Game 表示 Battlefield 游戏。
type Game int

// 支持的游戏。
const (
	GameBF1 Game = iota
	GameBFV
	GameBF2042
	GameBF6
	GameBF4
)

// String 返回游戏的配置名。
func (g Game) String() string {
	switch g {
	case GameBF1:
		return "bf1"
	case GameBFV:
		return "bfv"
	case GameBF2042:
		return "bf2042"
	case GameBF6:
		return "bf6"
	case GameBF4:
		return "bf4"
	default:
		return "bf1"
	}
}

// ParseGame 解析游戏名字符串（bf1 / bfv / bf4 / bf2042 / bf6）。
func ParseGame(s string) (Game, error) {
	switch s {
	case "bf1":
		return GameBF1, nil
	case "bfv":
		return GameBFV, nil
	case "bf4":
		return GameBF4, nil
	case "bf2042":
		return GameBF2042, nil
	case "bf6":
		return GameBF6, nil
	default:
		return GameBF1, fmt.Errorf("未知游戏: %s", s)
	}
}

// 武器类别（BF1 Progression 接口的 categoryId）。
const (
	WeaponALL     = "ALL"
	WeaponElite   = "ID_P_CAT_FIELDKIT"         // 精英兵
	WeaponLMG     = "ID_P_CAT_LMG"              // 轻机枪
	WeaponMelee   = "ID_P_CAT_MELEE"            // 近战武器
	WeaponGadget  = "ID_P_CAT_GADGET"           // 配备
	WeaponSemi    = "ID_P_CAT_SEMI"             // 半自动
	WeaponGrenade = "ID_P_CAT_GRENADE"          // 手榴弹
	WeaponSIR     = "ID_P_CAT_SIR"              // 制式步枪
	WeaponShotgun = "ID_P_CAT_SHOTGUN"          // 霰弹枪
	WeaponDriver  = "ID_P_CAT_VEHICLEKITWEAPON" // 驾驶员
	WeaponSMG     = "ID_P_CAT_SMG"              // 冲锋枪
	WeaponSidearm = "ID_P_CAT_SIDEARM"          // 手枪
	WeaponBolt    = "ID_P_CAT_BOLT"             // 步枪
)
