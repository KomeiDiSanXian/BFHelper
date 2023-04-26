// Package bf1rsp 战地服务器操作请求结构体
package bf1rsp

import bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"

// Post 结构体
type Post struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Game       string `json:"game"`
		PersonaID  string `json:"personaID"`
		GameID     string `json:"gameID"`
		ServerID   string `json:"serverID"`
		PGid       string `json:"persistedGameID"`
		LevelIndex int    `json:"levelIndex"`
		Reason     string `json:"reason"`
	} `json:"params"`
	ID string `json:"id"`
}

// NewPostKick 踢出游戏
func NewPostKick(pid, gid, reason string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  bf1api.KICK,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaID\""
			GameID     string `json:"gameID"`
			ServerID   string `json:"serverID"`
			PGid       string `json:"persistedGameID"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:      bf1api.BF1,
			PersonaID: pid,
			GameID:    gid,
			Reason:    reason,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// NewPostBan 服务器封禁
func NewPostBan(pid, sid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  bf1api.ADDBAN,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaID\""
			GameID     string `json:"gameID"`
			ServerID   string `json:"serverID"`
			PGid       string `json:"persistedGameID"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:      bf1api.BF1,
			PersonaID: pid,
			ServerID:  sid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// NewPostRemoveBan 服务器解封
func NewPostRemoveBan(pid, sid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  bf1api.REMOVEBAN,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaID\""
			GameID     string `json:"gameID"`
			ServerID   string `json:"serverID"`
			PGid       string `json:"persistedGameID"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:      bf1api.BF1,
			PersonaID: pid,
			ServerID:  sid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// NewPostChangeMap 换图
func NewPostChangeMap(pgid string, index int) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  bf1api.MAPS,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaID\""
			GameID     string `json:"gameID"`
			ServerID   string `json:"serverID"`
			PGid       string `json:"persistedGameID"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:       bf1api.BF1,
			PGid:       pgid,
			LevelIndex: index,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// NewPostGetServerDetails 获取完整服务器信息
func NewPostGetServerDetails(gid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  bf1api.SERVERDETALS,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaID\""
			GameID     string `json:"gameID"`
			ServerID   string `json:"serverID"`
			PGid       string `json:"persistedGameID"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:   bf1api.BF1,
			GameID: gid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// NewPostGetServerInfo 获取服务器部分信息
func NewPostGetServerInfo(gid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  bf1api.SERVERINFO,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaID\""
			GameID     string `json:"gameID"`
			ServerID   string `json:"serverID"`
			PGid       string `json:"persistedGameID"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:   bf1api.BF1,
			GameID: gid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// NewPostRSPInfo 获取服务器rsp信息
func NewPostRSPInfo(sid string) *Post {
	return &Post{
		Jsonrpc: "2.0",
		Method:  bf1api.SERVERRSP,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaID\""
			GameID     string `json:"gameID"`
			ServerID   string `json:"serverID"`
			PGid       string `json:"persistedGameID"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:     bf1api.BF1,
			ServerID: sid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}
