package bf1model

import "gorm.io/gorm"

/*
	按照群号绑定服务器，一个群有多个服务器
	暂不使用
*/

//群绑定服务器
type Group struct {
	gorm.Model
	GroupID int64 `gorm:"primaryKey"`
	Servers []Server
	Admins  []Admin
	Banlists []Banlist
}

//服务器 表
type Server struct {
	gorm.Model
	Gameid      string //gid
	Serverid    string //sid
	PGid        string //also guid
	NameInGroup string //群内对该服务器起的别名
	ServerName  string //服务器名
	Owner       string //腐竹
	Marks       int    //收藏数
	Synopsis    string //简介
}

//服务器管理
type Admin struct {
	gorm.Model
	QQid int64
}

//ban 表
type Banlist struct {
	gorm.Model
	Players []Player
}

type ServerDB gorm.DB
