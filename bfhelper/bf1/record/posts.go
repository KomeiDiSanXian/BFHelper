// Package bf1record 战地相关战绩查询结构体
package bf1record

import rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"

// Post 获取数据
type Post struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Param  `json:"params"`
	ID      string `json:"id"`
}

// Param post parameters
type Param struct {
	Game       string   `json:"game"`
	PersonaID  string   `json:"personaId"`
	PersonaIds []string `json:"personaIds"`
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
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
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
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
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
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
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
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
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
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}
