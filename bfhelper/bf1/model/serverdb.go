package bf1model

import (
	"errors"
	"time"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/rsp"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 群绑定服务器
type Group struct {
	GroupID   int64    `gorm:"primaryKey"`
	Servers   []Server `gorm:"foreignkey:GroupID;references:GroupID"`
	Admins    []Admin  `gorm:"foreignkey:GroupID;references:GroupID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 服务器 表
type Server struct {
	GroupID     int64
	Gameid      string `gorm:"primaryKey"` //gid
	Serverid    string `gorm:"primaryKey"` //sid
	PGid        string //also guid
	NameInGroup string //群内对该服务器起的别名
	ServerName  string //服务器名
	Owner       string //腐竹
	Bans        []Ban  `gorm:"foreignkey:Gameid;references:Gameid"` //Ban 列
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// 服务器管理
type Admin struct {
	GroupID   int64 `gorm:"primaryKey"`
	QQid      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BAN列 表
type Ban struct {
	Gameid      string `gorm:"primaryKey"`
	DisplayName string
	PersonalID  string `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ServerDB gorm.DB

// curd
// Create new server bind
func (sdb *ServerDB) Create(groupid int64, gameid string) error {
	// check gameid
	post := bf1rsp.NewPostGetServerDetails(gameid)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "error.code").String() == bf1api.ErrServerNotFound {
		return errors.New("找不到服务器，请检查gid")
	}
	// put in database
	// read needed info from data
	result := gjson.GetMany(
		data,
		"result.rspInfo.bannedList",
		"result.rspInfo.server.serverId",
		"result.rspInfo.server.persistedGameId",
		"result.rspInfo.server.name",
		"result.rspInfo.server.ownerId",
	)
	ban := result[0].Array()
	var bans []Ban
	for _, v := range ban {
		bans = append(bans, Ban{
			DisplayName: v.Get("displayName").Str,
			PersonalID:  v.Get("personaId").Str,
		})
	}
	grp := &Group{
		GroupID: groupid,
		Servers: []Server{
			{
				Gameid:     gameid,
				Serverid:   result[1].String(),
				PGid:       result[2].Str,
				ServerName: result[3].Str,
				Owner:      result[4].String(),
				Bans:       bans,
			},
		},
	}
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Create(grp).Error
}

// update
func (sdb *ServerDB) Update(grp Group) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Session(&gorm.Session{FullSaveAssociations: true}).Updates(&grp).Error
}

// read
func (sdb *ServerDB) Find(grpid int64) (*Group, error) {
	var result Group
	rmu.Lock()
	defer rmu.Unlock()
	err := (*gorm.DB)(sdb).Model(&Group{}).Where("group_id = ?", grpid).Preload("Servers.Bans").Preload(clause.Associations).First(&result).Error
	return &result, err
}

// add admin to group
func (sdb *ServerDB) AddAdmin(grpid, qid int64) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Where("group_id = ?", grpid).Association("Admins").Append(&Admin{QQid: qid})
}

// add server to group
func (sdb *ServerDB) AddServer(grpid int64, gid string) error {
	// check gameid
	post := bf1rsp.NewPostGetServerDetails(gid)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "error.code").String() == bf1api.ErrServerNotFound {
		return errors.New("找不到服务器，请检查gid")
	}
	// put in database
	// read needed info from data
	result := gjson.GetMany(
		data,
		"result.rspInfo.bannedList",
		"result.rspInfo.server.serverId",
		"result.rspInfo.server.persistedGameId",
		"result.rspInfo.server.name",
		"result.rspInfo.server.ownerId",
	)
	ban := result[0].Array()
	var bans []Ban
	for _, v := range ban {
		bans = append(bans, Ban{
			DisplayName: v.Get("displayName").Str,
			PersonalID:  v.Get("personaId").Str,
		})
	}
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Where("group_id = ?", grpid).Session(&gorm.Session{FullSaveAssociations: true}).Association("Servers").
		Append(&Server{
			Gameid:     gid,
			Serverid:   result[1].String(),
			PGid:       result[2].Str,
			ServerName: result[3].Str,
			Owner:      result[4].String(),
			Bans:       bans,
		})
}

// add ban to server
func (sdb *ServerDB) AddBan(gid, name string, grpid int64) error {
	pid, err := bf1api.GetPersonalID(name)
	if err != nil {
		return errors.New("添加到ban列失败：" + err.Error())
	}
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Where("group_id = ?", grpid).Session(&gorm.Session{FullSaveAssociations: true}).Association("Servers").
		Append(&Server{
			Gameid: gid,
			Bans: append([]Ban{}, Ban{
				Gameid:      gid,
				DisplayName: name,
				PersonalID:  pid,
			}),
		})
}

/*
// remove admin
func (sdb *ServerDB) DelAdmin(grpid, qid int64) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Where("group_id = ?", grpid).Association("Admins").Delete(&Admin{QQid: qid})
}

// remove server
func (sdb *ServerDB) DelServer(grpid int64, gid string) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Where("group_id = ?", grpid).Session(&gorm.Session{FullSaveAssociations: true}).Association("Servers").
		Delete(&Server{
			Gameid: gid,
		})
}

// remove ban
func (sdb *ServerDB) DelBan(gid, name string, grpid int64) error {
	rmu.Lock()
	defer rmu.Unlock()
	return (*gorm.DB)(sdb).Model(&Group{GroupID: grpid}).Where("group_id = ?", grpid).Session(&gorm.Session{FullSaveAssociations: true}).Association("Servers").
	Delete(&Server{
		Gameid: gid,
	})
}
*/
