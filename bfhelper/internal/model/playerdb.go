// Package model 玩家操作
package model

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

// Create 创建玩家条目
func (player *Player) Create(db *gorm.DB) error {
	return db.Create(player).Error
}

// Update 使用Updates方法更新玩家信息，0值不更新
func (player *Player) Update(db *gorm.DB) error {
	if player.Qid == 0 {
		return errors.New("qid cannot be empty")
	}

	return db.Model(&Player{}).Updates(player).Error
}

// Delete 删除玩家信息
func (player *Player) Delete(db *gorm.DB) error {
	if player.Qid == 0 {
		return errors.New("qid cannot be empty")
	}

	return db.Where("qid = ?", player.Qid).Delete(&Player{}).Error
}

// GetByQID 使用qq号查询玩家表
func (player *Player) GetByQID(db *gorm.DB, qid int64) (*Player, error) {
	if qid == 0 {
		return nil, errors.New("qid cannot be empty")
	}

	var playerResult Player
	if err := db.Where("qid = ?", qid).First(&playerResult).Error; err != nil {
		return nil, err
	}
	return &playerResult, nil
}

// GetByName 使用玩家名查询玩家表
func (player *Player) GetByName(db *gorm.DB, name string) (*Player, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	var playerResult Player
	if err := db.Where("display_name = ?", name).First(&playerResult).Error; err != nil {
		return nil, err
	}
	return &playerResult, nil
}
