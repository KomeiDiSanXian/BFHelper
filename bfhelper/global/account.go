// Package global 全局sesssion 信息
package global

type account struct {
	Username string
	Password string
	Session string // Session is X-Gatewaysession
	Token   string // Token is bearerAccessToken
	SID     string // SID is cookie sid
	Remid   string // Remid is cookie remid
}

// Sakura apikey 和api id
type Sakura struct {
	SakuraID    string
	SakuraToken string
}

// Account 全局账号登陆信息
//
// 存有 session, token, sid, remid
var Account *account

// SakuraAPI 含有api id, api key
var SakuraAPI *Sakura
