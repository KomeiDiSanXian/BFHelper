// Package battlefield 封装 GoBattlefieldAPI，提供 EA 登录与 Sparta 服务器管理能力。
package battlefield

import (
	"context"
	"fmt"
	"sync"

	"github.com/Dev4BF/GoBattlefieldAPI/auth"
	"github.com/Dev4BF/GoBattlefieldAPI/gameclient"
	"github.com/tidwall/gjson"
)

// Client 是 EA 认证 + Sparta 服务器管理的统一客户端。
type Client struct {
	mu      sync.Mutex
	session *auth.Session
	game    auth.Game
	gc      *gameclient.Client
}

// NewClient 创建客户端。
func NewClient() *Client { return &Client{} }

// Login 使用 sid/remid 完成 EA OAuth 登录。
//
// game 决定 Step 2/3 使用的游戏客户端配置；BF1/BFV/BF4 额外创建 Sparta 客户端。
func (c *Client) Login(ctx context.Context, sid, remid string, game auth.Game) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	session, err := auth.Login(sid, remid, game)
	if err != nil {
		return fmt.Errorf("EA 登录失败: %w", err)
	}

	var gc *gameclient.Client
	if game.IsSparta() {
		gc, err = newGameClient(game, session.GetGatewaySession())
		if err != nil {
			return fmt.Errorf("创建 Sparta 客户端失败: %w", err)
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.session = session
	c.game = game
	c.gc = gc
	return nil
}

// Logout 登出并清空会话。
func (c *Client) Logout() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.gc != nil {
		_, _ = c.gc.Logout()
	}
	c.session = nil
	c.game = 0
	c.gc = nil
}

// IsLoggedIn 返回是否已登录。
func (c *Client) IsLoggedIn() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.session != nil
}

// GatewaySession 返回当前 gateway session。
func (c *Client) GatewaySession() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.session == nil {
		return ""
	}
	return c.session.GetGatewaySession()
}

// Game 返回当前登录的游戏。
func (c *Client) Game() auth.Game {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.game
}

// LookupPersonaByName 按玩家名查询 Persona。
func (c *Client) LookupPersonaByName(ctx context.Context, name string) (auth.Persona, error) {
	if err := ctx.Err(); err != nil {
		return auth.Persona{}, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.session == nil {
		return auth.Persona{}, fmt.Errorf("未登录 EA 账号")
	}

	personas, err := c.session.LookupPersonaByNames(name)
	if err != nil {
		return auth.Persona{}, err
	}
	for _, p := range personas {
		if p.PersonaID != 0 {
			return p, nil
		}
	}
	if len(personas) > 0 {
		return personas[0], nil
	}
	return auth.Persona{}, fmt.Errorf("未找到玩家: %s", name)
}

// requireGC 返回 Sparta 客户端，未登录或游戏不支持时返回错误。
func (c *Client) requireGC() (*gameclient.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.gc == nil {
		return nil, fmt.Errorf("需要先登录 BF1/BFV/BF4 账号")
	}
	return c.gc, nil
}

// KickPlayer 踢出玩家（gameID 为游戏内服务器 ID）。
func (c *Client) KickPlayer(gameID, personaID, reason string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.KickPlayer(gameID, personaID, reason)
}

// BanPlayer 封禁玩家。
//
// gameID 为游戏内服务器 ID，内部解析 RSP serverId 后调用 addServerBan。
func (c *Client) BanPlayer(gameID, personaID string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	serverID, err := c.resolveServerID(gameID)
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.AddServerBan(serverID, personaID)
}

// UnbanPlayer 解封玩家。
func (c *Client) UnbanPlayer(gameID, personaID string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	serverID, err := c.resolveServerID(gameID)
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.RemoveServerBan(serverID, personaID)
}

// MovePlayer 将玩家移动到指定队伍（teamID 为 1 或 2）。
func (c *Client) MovePlayer(gameID, personaID string, teamID int) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.MovePlayer(gameID, personaID, teamID)
}

// AddServerAdmin 添加服务器管理员。
func (c *Client) AddServerAdmin(gameID, personaID string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	serverID, err := c.resolveServerID(gameID)
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.AddServerAdmin(serverID, personaID)
}

// RemoveServerAdmin 移除服务器管理员。
func (c *Client) RemoveServerAdmin(gameID, personaID string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	serverID, err := c.resolveServerID(gameID)
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.RemoveServerAdmin(serverID, personaID)
}

// AddServerVIP 添加服务器 VIP。
func (c *Client) AddServerVIP(gameID, personaID string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	serverID, err := c.resolveServerID(gameID)
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.AddServerVIP(serverID, personaID)
}

// RemoveServerVIP 移除服务器 VIP。
func (c *Client) RemoveServerVIP(gameID, personaID string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	serverID, err := c.resolveServerID(gameID)
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.RemoveServerVIP(serverID, personaID)
}

// SwitchMap 切换服务器地图（persistedGameID 为服务器 guid）。
func (c *Client) SwitchMap(persistedGameID string, levelIndex int) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.SwitchMap(persistedGameID, levelIndex)
}

// GetLogs 获取服务器操作日志。
func (c *Client) GetLogs(gameID string, offset, limit int) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	serverID, err := c.resolveServerID(gameID)
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.GetLogs(serverID, offset, limit)
}

// GetServerDetails 获取服务器完整信息（含地图轮换与 RSP 信息）。
func (c *Client) GetServerDetails(gameID string) (gjson.Result, error) {
	gc, err := c.requireGC()
	if err != nil {
		return gjson.Result{}, err
	}
	return gc.GetFullServerDetails(gameID)
}

// ServerInfo 服务器基本信息。
type ServerInfo struct {
	GameID     string
	ServerID   string
	PGID       string
	ServerName string
}

// GetServerInfo 获取服务器完整信息，用于群组绑定服务器。
func (c *Client) GetServerInfo(ctx context.Context, gameID string) (*ServerInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	result, err := c.GetServerDetails(gameID)
	if err != nil {
		return nil, err
	}
	info := &ServerInfo{
		ServerID:   result.Get("rspInfo.server.serverId").Str,
		GameID:     result.Get("serverInfo.gameId").Str,
		PGID:       result.Get("serverInfo.guid").Str,
		ServerName: result.Get("serverInfo.name").Str,
	}
	if info.GameID == "" {
		info.GameID = gameID
	}
	return info, nil
}

// resolveServerID 通过 gameId 解析 RSP serverId。
func (c *Client) resolveServerID(gameID string) (string, error) {
	result, err := c.GetServerDetails(gameID)
	if err != nil {
		return "", err
	}
	serverID := result.Get("rspInfo.server.serverId").Str
	if serverID == "" {
		return "", fmt.Errorf("服务器 %s 未开放 RSP 管理", gameID)
	}
	return serverID, nil
}

// newGameClient 按游戏创建 Sparta 客户端。
func newGameClient(game auth.Game, session string) (*gameclient.Client, error) {
	var cfg gameclient.Config
	switch game {
	case auth.GameBF1:
		cfg = gameclient.BF1Config()
	case auth.GameBFV:
		cfg = gameclient.BFVConfig()
	case auth.GameBF4:
		cfg = gameclient.BF4Config()
	case auth.GameBF2042, auth.GameBF6:
		// BF2042 / BF6 使用 gRPC-Web 网关，无 Sparta 客户端。
		cfg = gameclient.BF1Config()
	}
	return gameclient.New(cfg, session)
}
