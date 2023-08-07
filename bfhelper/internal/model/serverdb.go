// Package model 服务器数据操作
package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Group 群组表
type Group struct {
	GroupID   int64 `gorm:"primary_key;auto_increment:false"`
	Owner     int64
	Servers   []Server `gorm:"many2many:group_servers"`
	Admins    []Admin  `gorm:"many2many:group_admins"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Server 服务器表
type Server struct {
	GameID      string `gorm:"primary_key"`
	ServerID    string
	PGID        string
	NameInGroup string
	ServerName  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Admin 管理表
type Admin struct {
	QQID      int64 `gorm:"primary_key;auto_increment:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewGroup creates a new Group
func NewGroup(groupID int64) *Group {
	return &Group{GroupID: groupID}
}

// NewServer creates a new Server
func NewServer(gameID string) *Server {
	return &Server{GameID: gameID}
}

// NewAdmin creates a new Admin
func NewAdmin(qid int64) *Admin {
	return &Admin{QQID: qid}
}

// Create 创建一个新的 Group 记录
func (g *Group) Create(db *gorm.DB) error {
	if g.GroupID == 0 {
		return errors.New("invalid Group ID")
	}
	return db.Create(g).Error
}

// GetByID 根据 GroupID 获取 Group 记录
func (g *Group) GetByID(db *gorm.DB) (*Group, error) {
	if g.GroupID == 0 {
		return nil, errors.New("invalid Group ID")
	}
	err := db.Preload("Servers").Preload("Admins").First(&g, g.GroupID).Error
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Update 更新 Group 记录
//
// 注意不能修改为各类型零值
func (g *Group) Update(db *gorm.DB) error {
	if g.GroupID == 0 {
		return errors.New("invalid Group ID")
	}
	// 开始事务
	tx := db.Begin()
	if err := tx.Model(g).Updates(g).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 添加关联
	if err := tx.Model(g).Association("Servers").Append(g.Servers).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(g).Association("Admins").Append(g.Admins).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	return nil
}

// DeleteByID 根据 GroupID 删除 Group 记录
func (g *Group) DeleteByID(db *gorm.DB) error {
	if g.GroupID == 0 {
		return errors.New("invalid Group ID")
	}
	// 删除 Group 表中的记录
	if err := db.Delete(&Group{GroupID: g.GroupID}).Error; err != nil {
		return err
	}

	// 删除关联表中的数据
	if err := db.Exec("DELETE FROM group_servers WHERE group_group_id = ?", g.GroupID).Error; err != nil {
		return err
	}

	return db.Exec("DELETE FROM group_admins WHERE group_group_id = ?", g.GroupID).Error
}

// IsAdmin 检查是否为该服务器群管理
func (g *Group) IsAdmin(db *gorm.DB, qid int64) bool {
	grpdb, err := g.GetByID(db)
	if err != nil {
		return false
	}
	for _, admin := range grpdb.Admins {
		if qid == admin.QQID {
			return true
		}
	}
	return qid == grpdb.Owner
}

// IsOwner 检查是否为服务器拥有者
func (g *Group) IsOwner(db *gorm.DB, qid int64) bool {
	grpdb, err := g.GetByID(db)
	if err != nil {
		return false
	}
	return grpdb.Owner == qid
}

// Create 创建一个新的 Server 记录
func (s *Server) Create(db *gorm.DB) error {
	if s.GameID == "" {
		return errors.New("invalid GameID")
	}
	return db.Create(s).Error
}

// GetByGameID 根据 GameID 获取 Server 记录
func (s *Server) GetByGameID(db *gorm.DB) (*Server, error) {
	if s.GameID == "" {
		return nil, errors.New("invalid GameID")
	}
	err := db.First(&s, s.GameID).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Update 更新 Server 记录
func (s *Server) Update(db *gorm.DB) error {
	if s.GameID == "" {
		return errors.New("invalid GameID")
	}
	return db.Model(&Player{}).Updates(s).Error
}

// DeleteByGameID 根据 GameID 删除 Server 记录
func (s *Server) DeleteByGameID(db *gorm.DB) error {
	if s.GameID == "" {
		return errors.New("invalid GameID")
	}
	return db.Delete(&Server{GameID: s.GameID}).Error
}

// Create 创建一个新的 Admin 记录
func (a *Admin) Create(db *gorm.DB) error {
	if a.QQID == 0 {
		return errors.New("invalid QQ")
	}
	return db.Create(a).Error
}

// GetByQQID 根据 QQID 获取 Admin 记录
func (a *Admin) GetByQQID(db *gorm.DB) (*Admin, error) {
	if a.QQID == 0 {
		return nil, errors.New("invalid QQ")
	}
	err := db.First(a, a.QQID).Error
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Update 更新 Admin 记录
func (a *Admin) Update(db *gorm.DB) error {
	if a.QQID == 0 {
		return errors.New("invalid QQ")
	}
	return db.Model(&Admin{}).Updates(a).Error
}

// DeleteByQQID 根据 QQID 删除 Admin 记录
func (a *Admin) DeleteByQQID(db *gorm.DB) error {
	if a.QQID == 0 {
		return errors.New("invalid QQ")
	}
	return db.Delete(&Admin{QQID: a.QQID}).Error
}
