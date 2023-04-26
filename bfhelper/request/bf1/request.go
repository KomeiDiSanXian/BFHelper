// Package request 战地服务器操作请求
package request

import (
	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/uuid"
)

// Post 结构体
type Post struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Param  `json:"params"`
	ID      string `json:"id"`
}

// Param parameters
type Param struct {
	Game       string   `json:"game"`
	PersonaID  string   `json:"personaID"`
	PersonaIds []string `json:"personaIds"`
	GameID     string   `json:"gameID"`
	ServerID   string   `json:"serverID"`
	PGid       string   `json:"persistedGameID"`
	LevelIndex int      `json:"levelIndex"`
	Reason     string   `json:"reason"`
}

// NewPostKick 踢出游戏
func NewPostKick(pid, gid, reason string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.KICK,
		Params: Param{
			Game:      rsp.BF1,
			PersonaID: pid,
			GameID:    gid,
			Reason:    reason,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostBan 服务器封禁
func NewPostBan(pid, sid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.ADDBAN,
		Params: Param{
			Game:      rsp.BF1,
			PersonaID: pid,
			ServerID:  sid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostRemoveBan 服务器解封
func NewPostRemoveBan(pid, sid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.REMOVEBAN,
		Params: Param{
			Game:      rsp.BF1,
			PersonaID: pid,
			ServerID:  sid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostChangeMap 换图
func NewPostChangeMap(pgid string, index int) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.MAPS,
		Params: Param{
			Game:       rsp.BF1,
			PGid:       pgid,
			LevelIndex: index,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostGetServerDetails 获取完整服务器信息
func NewPostGetServerDetails(gid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.SERVERDETALS,
		Params: Param{
			Game:   rsp.BF1,
			GameID: gid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostGetServerInfo 获取服务器部分信息
func NewPostGetServerInfo(gid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.SERVERINFO,
		Params: Param{
			Game:   rsp.BF1,
			GameID: gid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostRSPInfo 获取服务器rsp信息
func NewPostRSPInfo(sid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.SERVERRSP,
		Params: Param{
			Game:     rsp.BF1,
			ServerID: sid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostWeapon 武器结构体
func NewPostWeapon(pid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.WEAPONS,
		Params: Param{
			Game:      rsp.BF1,
			PersonaID: pid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostVehicle 载具结构体
func NewPostVehicle(pid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.VEHICLES,
		Params: Param{
			Game:      rsp.BF1,
			PersonaID: pid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostRecent 最近游玩的服务器
func NewPostRecent(pid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.RECENTSERVER,
		Params: Param{
			Game:      rsp.BF1,
			PersonaID: pid,
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostPlaying 正在游玩
func NewPostPlaying(pid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.PLAYING,
		Params: Param{
			Game:       rsp.BF1,
			PersonaIds: []string{pid},
		},
		ID: uuid.NewUUID(),
	}
}

// NewPostStats 战绩
func NewPostStats(pid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  rsp.STATS,
		Params: Param{
			Game:      rsp.BF1,
			PersonaID: pid,
		},
		ID: uuid.NewUUID(),
	}
}
