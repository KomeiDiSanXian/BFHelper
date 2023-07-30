// Package bf1model bf1玩家操作
package bf1model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Player 玩家表
type Player struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PersonalID  string // pid
	Qid         int64  `gorm:"primary_key;auto_increment:false"` // QQ号
	DisplayName string // 玩家id
}

// PlayerRepository 玩家表数据库
type PlayerRepository struct {
	db *gorm.DB
}

// NewPlayerRepository 新建数据操作
func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

// Create 创建玩家条目
func (r *PlayerRepository) Create(player *Player) error {
	return r.db.Create(player).Error
}

// Update 使用Updates方法更新玩家信息，0值不更新
func (r *PlayerRepository) Update(player *Player) error {
	if player.Qid == 0 {
		return errors.New("qid cannot be empty")
	}

	return r.db.Model(&Player{}).Updates(player).Error
}

// Delete 删除玩家信息
func (r *PlayerRepository) Delete(qid int64) error {
	if qid == 0 {
		return errors.New("qid cannot be empty")
	}

	return r.db.Where("qid = ?", qid).Delete(&Player{}).Error
}

// GetByQID 使用qq号查询玩家表
func (r *PlayerRepository) GetByQID(qid int64) (*Player, error) {
	if qid == 0 {
		return nil, errors.New("qid cannot be empty")
	}

	var player Player
	if err := r.db.Where("qid = ?", qid).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

// GetByName 使用玩家名查询玩家表
func (r *PlayerRepository) GetByName(name string) (*Player, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	var player Player
	if err := r.db.Where("display_name = ?", name).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}
