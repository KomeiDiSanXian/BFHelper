package bf1rsp

import bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"

type post struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Game       string `json:"game"`
		PersonaID  string `json:"personaId"`
		GameId     string `json:"gameId"`
		ServerId   string `json:"serverId"`
		PGid       string `json:"persistedGameId"`
		LevelIndex int    `json:"levelIndex"`
		Reason     string `json:"reason"`
	} `json:"params"`
	ID string `json:"id"`
}

// POST 踢出游戏
func NewPostKick(pid, gid, reason string) *post {
	return &post{
		Jsonrpc: "2.0",
		Method:  bf1api.KICK,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaId\""
			GameId     string `json:"gameId"`
			ServerId   string `json:"serverId"`
			PGid       string `json:"persistedGameId"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:      bf1api.BF1,
			PersonaID: pid,
			GameId:    gid,
			Reason:    reason,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// POST 服务器封禁
func NewPostBan(pid, sid string) *post {
	return &post{
		Jsonrpc: "2.0",
		Method:  bf1api.ADDBAN,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaId\""
			GameId     string `json:"gameId"`
			ServerId   string `json:"serverId"`
			PGid       string `json:"persistedGameId"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:      bf1api.BF1,
			PersonaID: pid,
			ServerId:  sid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// POST 服务器解封
func NewPostRemoveBan(pid, sid string) *post {
	return &post{
		Jsonrpc: "2.0",
		Method:  bf1api.REMOVEBAN,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaId\""
			GameId     string `json:"gameId"`
			ServerId   string `json:"serverId"`
			PGid       string `json:"persistedGameId"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:      bf1api.BF1,
			PersonaID: pid,
			ServerId:  sid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// POST 换图
func NewPostChangeMap(pgid string, index int) *post {
	return &post{
		Jsonrpc: "2.0",
		Method:  bf1api.MAPS,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaId\""
			GameId     string `json:"gameId"`
			ServerId   string `json:"serverId"`
			PGid       string `json:"persistedGameId"`
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

// POST 获取完整服务器信息
func NewPostGetServerDetails(gid string) *post {
	return &post{
		Jsonrpc: "2.0",
		Method:  bf1api.SERVERDETALS,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaId\""
			GameId     string `json:"gameId"`
			ServerId   string `json:"serverId"`
			PGid       string `json:"persistedGameId"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:   bf1api.BF1,
			GameId: gid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}

// POST 获取服务器部分信息
func NewPostGetServerInfo(gid string) *post {
	return &post{
		Jsonrpc: "2.0",
		Method:  bf1api.SERVERINFO,
		Params: struct {
			Game       string "json:\"game\""
			PersonaID  string "json:\"personaId\""
			GameId     string `json:"gameId"`
			ServerId   string `json:"serverId"`
			PGid       string `json:"persistedGameId"`
			LevelIndex int    `json:"levelIndex"`
			Reason     string `json:"reason"`
		}{
			Game:   bf1api.BF1,
			GameId: gid,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
}
