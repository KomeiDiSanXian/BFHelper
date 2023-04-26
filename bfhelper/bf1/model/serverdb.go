// Package bf1model 战地服务器相关数据库操作
package bf1model

import (
	"time"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/rsp"
)

// Group 群绑定服务器
type Group struct {
	GroupID   int64    `gorm:"primaryKey"`
	Owner     int64    `gorm:"not null"`
	Servers   []Server `gorm:"foreignkey:GroupID;references:GroupID"`
	Admins    []Admin  `gorm:"foreignkey:GroupID;references:GroupID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Server 表
type Server struct {
	GroupID     int64  `gorm:"primaryKey"`
	Gameid      string `gorm:"primaryKey"` // gid
	Serverid    string // sid
	PGid        string // also guid
	NameInGroup string // 群内对该服务器起的别名
	ServerName  string // 服务器名
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Admin 服务器管理
type Admin struct {
	GroupID   int64
	QQid      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ServerDB 服务器数据
type ServerDB gorm.DB

// Create new server bind
func (sdb *ServerDB) Create(groupid, ownerid int64, gameid string) error {
	// check gameid
	post := bf1rsp.NewPostGetServerDetails(gameid)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	err = bf1api.Exception(gjson.Get(data, "error.code").Int())
	if err != nil {
		return err
	}
	// put in database
	// read needed info from data
	result := gjson.GetMany(
		data,
		"result.rspInfo.server.serverId",
		"result.rspInfo.server.persistedGameId",
		"result.rspInfo.server.name",
	)
	grp := &Group{
		GroupID: groupid,
		Owner:   ownerid,
		Servers: []Server{
			{
				Gameid:     gameid,
				Serverid:   result[0].Str,
				PGid:       result[1].Str,
				ServerName: result[2].Str,
			},
		},
	}
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Create(grp).Error
}

// Update 更新群组数据
func (sdb *ServerDB) Update(grp Group) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Session(&gorm.Session{FullSaveAssociations: true}).Updates(&grp).Error
}

// Find 寻找群组
func (sdb *ServerDB) Find(grpid int64) (*Group, error) {
	var result Group
	rmu.Lock()
	defer rmu.Unlock()
	err := (*gorm.DB)(sdb).Model(&Group{}).Where("group_id = ?", grpid).Preload(clause.Associations).First(&result).Error
	return &result, err
}

// AddAdmin 添加管理员到指定群
func (sdb *ServerDB) AddAdmin(grpid, qid int64) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Association("Admins").Append(&Admin{QQid: qid})
}

// AddServer 添加服务器到指定群
func (sdb *ServerDB) AddServer(grpid int64, gid string) error {
	// check gameid
	post := bf1rsp.NewPostGetServerDetails(gid)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	err = bf1api.Exception(gjson.Get(data, "error.code").Int())
	if err != nil {
		return err
	}
	// put in database
	// read needed info from data
	result := gjson.GetMany(
		data,
		"result.rspInfo.server.serverId",
		"result.rspInfo.server.persistedGameId",
		"result.rspInfo.server.name",
	)
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Association("Servers").
		Append(&Server{
			Gameid:     gid,
			Serverid:   result[0].Str,
			PGid:       result[1].Str,
			ServerName: result[2].Str,
		})
}

// SetAlias 设置服务器别名
func (sdb *ServerDB) SetAlias(grpid int64, gid, alias string) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Where("group_id = ? AND gameid = ?", grpid, gid).Updates(&Server{NameInGroup: alias}).Error
}

// GetServer 由别名获取服务器信息
func (sdb *ServerDB) GetServer(alias string, grpid int64) (*Server, error) {
	var s Server
	rmu.Lock()
	defer rmu.Unlock()
	err := (*gorm.DB)(sdb).Where("group_id = ? AND name_in_group = ?", grpid, alias).First(&s).Error
	return &s, err
}

// DelAdmin 删除群管理
func (sdb *ServerDB) DelAdmin(grpid, qid int64) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Where("group_id = ? AND q_qid = ?", grpid, qid).Delete(&Admin{}).Error
}

// DelServer 删除群服务器
func (sdb *ServerDB) DelServer(grpid int64, gid string) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Where("group_id = ? AND gameid = ?", grpid, gid).Delete(&Server{}).Error
}

// ChangeOwner 更改群拥有人
func (sdb *ServerDB) ChangeOwner(grpid, owner int64) error {
	return sdb.Update(Group{GroupID: grpid, Owner: owner})
}

// Close 关闭数据库连接
func (sdb *ServerDB) Close() error {
	sqlDB, err := (*gorm.DB)(sdb).DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// IsAdmin 是否为群管理
func (sdb *ServerDB) IsAdmin(grpid, qid int64) (bool, error) {
	d, err := sdb.Find(grpid)
	if err != nil {
		return false, err
	}
	for _, v := range d.Admins {
		if v.QQid == qid {
			return true, nil
		}
		continue
	}
	return d.Owner == qid, nil
}

// IsOwner 是否为群主
func (sdb *ServerDB) IsOwner(grpid, qid int64) (bool, error) {
	d, err := sdb.Find(grpid)
	if err != nil {
		return false, err
	}
	return d.Owner == qid, nil
}
