package core

import "errors"

// 通用错误哨兵。
var (
	// ErrPlayerNotFound 玩家未绑定或不存在。
	ErrPlayerNotFound = errors.New("未绑定账号，请先使用 /bind <玩家名> 绑定")
)
