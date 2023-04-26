// Package bf1model 战地玩家相关数据库操作
package bf1model

import "gorm.io/gorm"

// Player 玩家表
type Player struct {
	gorm.Model
	PersonalID  string `gorm:"primaryKey"` // pid
	Qid         int64  `gorm:"primaryKey"` // QQ号
	DisplayName string // 玩家id
	IsHack      bool   // 是否实锤开挂
	BanReason   string
}

// PlayerDB 玩家数据库
type PlayerDB gorm.DB

// CURD...

// Create 创建玩家数据
func (pdb *PlayerDB) Create(player Player) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(pdb).Create(&player).Error
}

// Update 根据qid来更新数据
func (pdb *PlayerDB) Update(player Player) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(pdb).Model(&Player{}).Where("qid = ?", player.Qid).Updates(&player).Error
}

// FindByPid 根据pid寻找玩家
func (pdb *PlayerDB) FindByPid(pid uint) (*Player, error) {
	var player Player
	rmu.Lock()
	defer rmu.Unlock()
	err := (*gorm.DB)(pdb).Model(&Player{}).First(&player, "personal_id = ?", pid).Error
	return &player, err
}

// FindByName 根据name寻找玩家
func (pdb *PlayerDB) FindByName(name string) (*Player, error) {
	var player Player
	rmu.Lock()
	defer rmu.Unlock()
	err := (*gorm.DB)(pdb).Model(&Player{}).First(&player, "display_name = ?", name).Error
	return &player, err
}

// FindByQid 根据qid寻找玩家
func (pdb *PlayerDB) FindByQid(qid int64) (*Player, error) {
	var player Player
	rmu.Lock()
	defer rmu.Unlock()
	err := (*gorm.DB)(pdb).Model(&Player{}).First(&player, "qid = ?", qid).Error
	return &player, err
}

// Delete 由pid删除玩家数据
func (pdb *PlayerDB) Delete(pid uint) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(pdb).Delete(&Player{}, "personal_id = ?", pid).Error
}

// Close the database
func (pdb *PlayerDB) Close() error {
	sqlDB, err := (*gorm.DB)(pdb).DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
