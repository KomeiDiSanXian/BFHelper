// Package model 服务器数据操作
//
// TODO: refactor
package model

import (
	"time"

	"github.com/jinzhu/gorm"
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

// GroupRepository 定义 Group 表的存储库
type GroupRepository struct {
	db *gorm.DB
}

// NewGroupRepository 创建 GroupRepository 实例
func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// CreateGroup 创建一个新的 Group 记录
func (r *GroupRepository) CreateGroup(group *Group) error {
	return r.db.Create(group).Error
}

// GetGroupByID 根据 GroupID 获取 Group 记录
func (r *GroupRepository) GetGroupByID(groupID int64) (*Group, error) {
	var group Group
	err := r.db.Preload("Servers").Preload("Admins").First(&group, groupID).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateGroup 更新 Group 记录
//
// 注意不能修改为各类型零值
func (r *GroupRepository) UpdateGroup(group *Group) error {
	// 开始事务
	tx := r.db.Begin()
	if err := tx.Model(group).Updates(group).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 添加关联
	if err := tx.Model(group).Association("Servers").Append(group.Servers).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(group).Association("Admins").Append(group.Admins).Error; err != nil {
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

// DeleteGroupByID 根据 GroupID 删除 Group 记录
func (r *GroupRepository) DeleteGroupByID(groupID int64) error {
	// 删除 Group 表中的记录
	if err := r.db.Delete(&Group{GroupID: groupID}).Error; err != nil {
		return err
	}

	// 删除关联表中的数据
	if err := r.db.Exec("DELETE FROM group_servers WHERE group_group_id = ?", groupID).Error; err != nil {
		return err
	}

	return r.db.Exec("DELETE FROM group_admins WHERE group_group_id = ?", groupID).Error
}

// IsGroupAdmin 检查是否为该服务器群管理
func (r *GroupRepository) IsGroupAdmin(groupID, qid int64) bool {
	grpdb, err := r.GetGroupByID(groupID)
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

// IsGroupOwner 检查是否为服务器拥有者
func (r *GroupRepository) IsGroupOwner(groupID, qid int64) bool {
	grpdb, err := r.GetGroupByID(groupID)
	if err != nil {
		return false
	}
	return grpdb.Owner == qid
}

// ServerRepository 定义 Server 表的存储库
type ServerRepository struct {
	db *gorm.DB
}

// NewServerRepository 创建 ServerRepository 实例
func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

// CreateServer 创建一个新的 Server 记录
func (r *ServerRepository) CreateServer(server *Server) error {
	return r.db.Create(server).Error
}

// GetServerByGameID 根据 GameID 获取 Server 记录
func (r *ServerRepository) GetServerByGameID(gameID string) (*Server, error) {
	var server Server
	err := r.db.First(&server, gameID).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// UpdateServer 更新 Server 记录
func (r *ServerRepository) UpdateServer(server *Server) error {
	return r.db.Save(server).Error
}

// DeleteServerByGameID 根据 GameID 删除 Server 记录
func (r *ServerRepository) DeleteServerByGameID(gameID string) error {
	return r.db.Delete(&Server{GameID: gameID}).Error
}

// AdminRepository 定义 Admin 表的存储库
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建 AdminRepository 实例
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// CreateAdmin 创建一个新的 Admin 记录
func (r *AdminRepository) CreateAdmin(admin *Admin) error {
	return r.db.Create(admin).Error
}

// GetAdminByQQID 根据 QQID 获取 Admin 记录
func (r *AdminRepository) GetAdminByQQID(qqid int64) (*Admin, error) {
	var admin Admin
	err := r.db.First(&admin, qqid).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateAdmin 更新 Admin 记录
func (r *AdminRepository) UpdateAdmin(admin *Admin) error {
	return r.db.Save(admin).Error
}

// DeleteAdminByQQID 根据 QQID 删除 Admin 记录
func (r *AdminRepository) DeleteAdminByQQID(qqid int64) error {
	return r.db.Delete(&Admin{QQID: qqid}).Error
}
