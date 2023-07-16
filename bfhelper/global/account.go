// Package global 全局sesssion 信息
package global

import "github.com/KomeiDiSanXian/BFHelper/bfhelper/setting"

type loginInfo struct {
	Session string // Session is X-Gatewaysession
	Token   string // Token is bearerAccessToken
	SID     string // SID is cookie sid
	Remid   string // Remid is cookie remid
}

type account struct {
	LoginedUser *setting.Account // LoginedUser 账号信息
	Info        *loginInfo       // Info 登陆信息
}

// Account 全局账号登陆信息
//
// 存有 session, token, sid, remid
var Account *account

// SakuraAPI
//
// 含有api id, api key
var SakuraAPI *setting.SakuraAPI