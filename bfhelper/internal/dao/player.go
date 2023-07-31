// Package dao 玩家dao层操作
package dao

import "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/model"

// CreatePlayer 创建玩家条目
func (d *Dao) CreatePlayer(qid int64, name string) error {
	player := model.Player{Qid: qid, DisplayName: name}
	return player.Create(d.engine)
}

// DeletePlayer 删除玩家
func (d *Dao) DeletePlayer(qid int64) error {
	player := model.Player{Qid: qid}
	return player.Delete(d.engine)
}

// UpdatePlayer 更新玩家信息
func (d *Dao) UpdatePlayer(qid int64, pid, name string) error {
	player := model.Player{Qid: qid, PersonalID: pid, DisplayName: name}
	return player.Update(d.engine)
}

// GetPlayerByQID 根据qq号获取玩家信息
func (d *Dao) GetPlayerByQID(qid int64) (*model.Player, error) {
	player := model.Player{Qid: qid}
	return player.GetByQID(d.engine, qid)
}

// GetPlayerByName 根据玩家名获取玩家信息
func (d *Dao) GetPlayerByName(name string) (*model.Player, error) {
	player := model.Player{DisplayName: name}
	return player.GetByName(d.engine, name)
}
