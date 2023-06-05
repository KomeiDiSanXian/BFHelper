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

type PlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

func (r *PlayerRepository) Create(player *Player) error {
	if err := r.db.Create(player).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) Update(player *Player) error {
	if player.Qid == 0 {
		return errors.New("qid cannot be empty")
	}

	if err := r.db.Model(&Player{}).Updates(player).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) Delete(qid int64) error {
	if qid == 0 {
		return errors.New("qid cannot be empty")
	}

	if err := r.db.Where("qid = ?", qid).Delete(&Player{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) GetByQID(qid int64) (*Player, error) {
	if qid == 0 {
		return nil, errors.New("qid cannot be empty")
	}

	var player Player
	if err := r.db.Where("qid = ?", qid).First(&player).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("Player not found")
		}
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepository) GetByName(name string) (*Player, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	var player Player
	if err := r.db.Where("display_name = ?", name).First(&player).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("Player not found")
		}
		return nil, err
	}
	return &player, nil
}
