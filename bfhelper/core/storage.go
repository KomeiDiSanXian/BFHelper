package core

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// PlayerStore 玩家数据的存储接口。
type PlayerStore interface {
	GetPlayerByQQ(qqID int64) (*Player, error)
	GetPlayerByName(name string) (*Player, error)
	CreatePlayer(p *Player) error
	DeletePlayer(qqID int64) error
}

// GroupStore 群组数据的存储接口。
type GroupStore interface {
	GetGroup(groupID int64) (*Group, error)
	CreateGroup(g *Group) error
	UpdateGroup(g *Group) error
	DeleteGroup(groupID int64) error
}

// ServerStore 服务器绑定数据的存储接口。
type ServerStore interface {
	AddServerToGroup(groupID int64, sv *Server) error
	RemoveServerFromGroup(groupID int64, gameID string) error
	GetServerByAlias(groupID int64, alias string) (*Server, error)
}

// AdminStore 群组管理员数据的存储接口。
type AdminStore interface {
	AddAdminToGroup(groupID int64, qqIDs ...int64) error
	RemoveAdminFromGroup(groupID int64, qqID int64) error
	IsServerOwner(groupID, qqID int64) bool
	IsServerAdmin(groupID, qqID int64) bool
}

// BlazeWatchStore 服务器监听数据的存储接口。
type BlazeWatchStore interface {
	SetBlazeWatch(groupID int64, gameID string, enabled bool) error
	RemoveBlazeWatch(groupID int64, gameID string) error
	GetBlazeWatches(groupID int64) ([]BlazeWatch, error)
}

// OperationLogStore 操作日志的存储接口。
type OperationLogStore interface {
	RecordOperationLog(log *OperationLog) error
	QueryOperationLogs(groupID int64, from, to *time.Time, action string, limit int) ([]OperationLog, error)
}

// DBProvider 暴露底层数据库句柄，供状态统计等需要直接访问数据库的场景使用。
type DBProvider interface {
	DB() *gorm.DB
}

// Storage 聚合全部存储能力，供需要完整数据访问权限的组件使用。
type Storage interface {
	PlayerStore
	GroupStore
	ServerStore
	AdminStore
	BlazeWatchStore
	OperationLogStore
	DBProvider
}

// NewStorage 基于 GORM 创建存储实现。
func NewStorage(db *gorm.DB) Storage {
	return &gormStorage{db: db}
}

type gormStorage struct {
	db *gorm.DB
}

func (s *gormStorage) DB() *gorm.DB { return s.db }

// ── Player ──────────────────────────────────────────────────────────────────

func (s *gormStorage) GetPlayerByQQ(qqID int64) (*Player, error) {
	var p Player
	if err := s.db.Where("qq_id = ?", qqID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *gormStorage) GetPlayerByName(name string) (*Player, error) {
	var p Player
	if err := s.db.Where("persona_name = ?", name).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *gormStorage) CreatePlayer(p *Player) error {
	if p == nil {
		return errors.New("nil player")
	}
	return s.db.Create(p).Error
}

func (s *gormStorage) DeletePlayer(qqID int64) error {
	return s.db.Where("qq_id = ?", qqID).Delete(&Player{}).Error
}

// ── Group ───────────────────────────────────────────────────────────────────

func (s *gormStorage) GetGroup(groupID int64) (*Group, error) {
	var g Group
	if err := s.db.Preload("Servers").Preload("Admins").Where("group_id = ?", groupID).First(&g).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *gormStorage) CreateGroup(g *Group) error {
	if g == nil {
		return errors.New("nil group")
	}
	return s.db.Create(g).Error
}

func (s *gormStorage) UpdateGroup(g *Group) error {
	if g == nil {
		return errors.New("nil group")
	}
	return s.db.Model(&Group{}).Where("group_id = ?", g.GroupID).Update("owner_qq", g.OwnerQQ).Error
}

func (s *gormStorage) DeleteGroup(groupID int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", groupID).Delete(&Server{}).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", groupID).Delete(&Admin{}).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", groupID).Delete(&BlazeWatch{}).Error; err != nil {
			return err
		}
		return tx.Where("group_id = ?", groupID).Delete(&Group{}).Error
	})
}

// ── Server ──────────────────────────────────────────────────────────────────

func (s *gormStorage) AddServerToGroup(groupID int64, sv *Server) error {
	if sv == nil {
		return errors.New("nil server")
	}
	sv.GroupID = groupID

	var existing Server
	err := s.db.Where("group_id = ? AND game_id = ?", groupID, sv.GameID).First(&existing).Error
	if err == nil {
		// 重复绑定视为更新，保留已有别名。
		sv.ID = existing.ID
		if sv.Alias == "" {
			sv.Alias = existing.Alias
		}
		if sv.ServerID == "" {
			sv.ServerID = existing.ServerID
		}
		return s.db.Save(sv).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return s.db.Create(sv).Error
}

func (s *gormStorage) RemoveServerFromGroup(groupID int64, gameID string) error {
	return s.db.Where("group_id = ? AND game_id = ?", groupID, gameID).Delete(&Server{}).Error
}

func (s *gormStorage) GetServerByAlias(groupID int64, alias string) (*Server, error) {
	var sv Server
	// 别名或 gameId 均可定位服务器。
	if err := s.db.Where("group_id = ? AND (alias = ? OR game_id = ?)", groupID, alias, alias).First(&sv).Error; err != nil {
		return nil, err
	}
	return &sv, nil
}

// ── Admin ───────────────────────────────────────────────────────────────────

func (s *gormStorage) AddAdminToGroup(groupID int64, qqIDs ...int64) error {
	if len(qqIDs) == 0 {
		return nil
	}
	admins := make([]Admin, 0, len(qqIDs))
	for _, qqID := range qqIDs {
		admins = append(admins, Admin{GroupID: groupID, QQID: qqID})
	}
	return s.db.Create(&admins).Error
}

func (s *gormStorage) RemoveAdminFromGroup(groupID int64, qqID int64) error {
	return s.db.Where("group_id = ? AND qq_id = ?", groupID, qqID).Delete(&Admin{}).Error
}

func (s *gormStorage) IsServerOwner(groupID, qqID int64) bool {
	var count int64
	s.db.Model(&Group{}).Where("group_id = ? AND owner_qq = ?", groupID, qqID).Count(&count)
	return count > 0
}

func (s *gormStorage) IsServerAdmin(groupID, qqID int64) bool {
	var count int64
	s.db.Model(&Admin{}).Where("group_id = ? AND qq_id = ?", groupID, qqID).Count(&count)
	return count > 0
}

// ── BlazeWatch ──────────────────────────────────────────────────────────────

func (s *gormStorage) SetBlazeWatch(groupID int64, gameID string, enabled bool) error {
	var w BlazeWatch
	err := s.db.Where("group_id = ? AND game_id = ?", groupID, gameID).First(&w).Error
	if err == nil {
		w.Enabled = enabled
		return s.db.Save(&w).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	w = BlazeWatch{GroupID: groupID, GameID: gameID, Enabled: enabled}
	return s.db.Create(&w).Error
}

func (s *gormStorage) RemoveBlazeWatch(groupID int64, gameID string) error {
	err := s.db.Where("group_id = ? AND game_id = ?", groupID, gameID).Delete(&BlazeWatch{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (s *gormStorage) GetBlazeWatches(groupID int64) ([]BlazeWatch, error) {
	var watches []BlazeWatch
	err := s.db.Where("group_id = ? AND enabled = ?", groupID, true).Find(&watches).Error
	return watches, err
}

// ── OperationLog ──────────────────────────────────────────────────────────────

func (s *gormStorage) RecordOperationLog(log *OperationLog) error {
	if log == nil {
		return errors.New("nil operation log")
	}
	return s.db.Create(log).Error
}

// QueryOperationLogs 按群组查询操作日志，支持时间范围（from/to，含边界）与操作类型过滤。
func (s *gormStorage) QueryOperationLogs(groupID int64, from, to *time.Time, action string, limit int) ([]OperationLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := s.db.Where("group_id = ?", groupID)
	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}
	if action != "" {
		q = q.Where("action = ?", action)
	}
	var logs []OperationLog
	err := q.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
