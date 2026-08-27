package core

import "time"

// Credentials EA 账号凭据。
type Credentials struct {
	SID   string `json:"sid"`
	Remid string `json:"remid"`
	Game  Game   `json:"game"`
}

// Player 绑定的玩家。
type Player struct {
	ID          uint  `gorm:"primaryKey"`
	QQID        int64 `gorm:"uniqueIndex"`
	PersonaName string
	PersonaID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Group 服务器群组。
type Group struct {
	ID      uint  `gorm:"primaryKey"`
	GroupID int64 `gorm:"uniqueIndex"`
	OwnerQQ int64
	Servers []Server `gorm:"foreignKey:GroupID"`
	Admins  []Admin  `gorm:"foreignKey:GroupID"`
	// 群内全部服务器
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Server 群组绑定的服务器。
type Server struct {
	ID      uint  `gorm:"primaryKey"`
	GroupID int64 `gorm:"index"`
	// GameID 游戏内服务器 ID（gameId）
	GameID string
	// ServerID RSP serverId（封禁/解封使用）
	ServerID string
	// PGID persistedGameId（切图使用）
	PGID string
	// ServerName 服务器名称
	ServerName string
	// Alias 群内别名
	Alias     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Admin 群组管理员。
type Admin struct {
	ID        uint  `gorm:"primaryKey"`
	GroupID   int64 `gorm:"index"`
	QQID      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BlazeWatch Blaze 实时监听项。
type BlazeWatch struct {
	ID        uint  `gorm:"primaryKey"`
	GroupID   int64 `gorm:"index"`
	GameID    string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OperationLog 群内管理操作日志。
type OperationLog struct {
	ID         uint  `gorm:"primaryKey"`
	GroupID    int64 `gorm:"index;not null"`
	OperatorQQ int64 `gorm:"index;not null"`
	// Action 操作类型：kick / ban / unban / banall / unbanall / cm / move /
	// sadmin_add / sadmin_rm / vip_add / vip_rm / group_create 等。
	Action string `gorm:"index;not null"`
	// Target 操作对象（玩家名 / QQ 号 / 服务器 GameID 等）。
	Target string `gorm:"index"`
	// ServerName 涉及的服务器名称（别名），多服操作时可能为空。
	ServerName string
	// Detail 补充说明（原因、结果等）。
	Detail  string
	Success bool
	// CreatedAt 操作时间。
	CreatedAt time.Time
}
