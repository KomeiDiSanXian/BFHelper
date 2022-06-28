package bf1record

import (
	"errors"

	"github.com/tidwall/gjson"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
)

//APIs
const (
	AuthAPI      string = "https://api.s-wg.net/ServersCollection/getPlayerAll?DisplayName="
	NativeAPI    string = "https://sparta-gw.battlelog.com/jsonrpc/pc/api"
	SessionAPI   string = "https://battlefield-api.sakurakooi.cyou/account/login"
	OperationAPI string = "https://sparta-gw.battlelog.com/jsonrpc/ps4/api" //交换和行动包查询
)

//ea账密
const (
	USERNAME string = "" //邮箱
	PASSWORD string = ""
)

//原生API 方法名常量
const (
	//NativeAPI
	ADDVIP    string = "RSP.addServerVip"
	REMOVEVIP string = "RSP.removeServerVip"
	ADDBAN    string = "RSP.addServerBan"
	REMOVEBAN string = "RSP.removeServerBan"
	KICK      string = "RSP.kickPlayer"
	STATS     string = "Stats.detailedStatsByPersonaId"
	WEAPONS   string = "Progression.getWeaponsByPersonaId"
	VEHICLES  string = "Progression.getVehiclesByPersonaId"
	//OperationAPI
	EXCHANGE string = "ScrapExchange.getOffers"
	CAMPAIGN string = "CampaignOperations.getPlayerCampaignStatus"
)

//设置session
var SESSION string = ""

var client = gentleman.New()

//Session 获取
func Session(username, password string, refreshToken bool) error {
	login := map[string]interface{}{"username": username, "password": password, "refreshToken": refreshToken}
	//requesting..
	client.URL(SessionAPI)
	client.Use(body.JSON(login))
	//寻找SakuraKooi申请APIKey...
	client.Use(headers.Set("Sakura-Instance-Id", ""))
	client.Use(headers.Set("Sakura-Access-Token", ""))

	res, err := client.Request().Method("POST").Send()
	if err != nil {
		return errors.New("更新session时出错：" + err.Error())
	}
	SESSION = gjson.Get(res.String(), "data.gatewaySession").Str
	return nil
}

//NativeAPI 返回json
func ReturnJson(url, method string, parms interface{}) (string, error) {
	client.URL(url)
	client.Use(body.JSON(parms))
	client.Use(headers.Set("X-Gatewaysession", SESSION))
	res, err := client.Request().Method(method).Send()
	data := res.String()
	code := gjson.Get(data, "error.code").Int()
	//如果session过期，重新请求
	if code == -32501 {
		Session(USERNAME, PASSWORD, true)
		return ReturnJson(url, method, parms)
	}
	if err != nil {
		return "", err
	}
	return data, nil
}
