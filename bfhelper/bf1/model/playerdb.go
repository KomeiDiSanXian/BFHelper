package bf1model

import "gorm.io/gorm"

//Player 表
type Player struct {
	gorm.Model
	PersonalID  string `gorm:"primaryKey"` //pid
	Qid         int64  `gorm:"primaryKey"` //QQ号
	DisplayName string //玩家id
	IsHack      bool   //是否实锤开挂
}

type PlayerDB gorm.DB

//CURD...

func (pdb *PlayerDB) Create(player Player) error {
	return (*gorm.DB)(pdb).Debug().Create(&player).Error
}

func (pdb *PlayerDB) Update(player Player) error {
	return (*gorm.DB)(pdb).Debug().Model(&Player{}).Where("qid = ?", player.Qid).Updates(&player).Error
}

func (pdb *PlayerDB) FindByPid(pid uint) (*Player, error) {
	var player Player
	err := (*gorm.DB)(pdb).Debug().Model(&Player{}).First(&player, "personal_id = ?", pid).Error
	return &player, err
}

func (pdb *PlayerDB) FindByName(name string) (*Player, error) {
	var player Player
	err := (*gorm.DB)(pdb).Debug().Model(&Player{}).First(&player, "display_name = ?", name).Error
	return &player, err
}

func (pdb *PlayerDB) FindByQid(qid int64) (*Player, error) {
	var player Player
	err := (*gorm.DB)(pdb).Debug().Model(&Player{}).First(&player, "qid = ?", qid).Error
	return &player, err
}

func (pdb *PlayerDB) Delete(pid uint) error {
	return (*gorm.DB)(pdb).Debug().Delete(&Player{}, "personal_id = ?", pid).Error
}
