package bf1api

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

// APIs
const (
	AuthAPI      string = "https://api.s-wg.net/ServersCollection/getPlayerAll?DisplayName="
	NativeAPI    string = "https://sparta-gw.battlelog.com/jsonrpc/pc/api"
	SessionAPI   string = "https://battlefield-api.sakurakooi.cyou/account/login"
	OperationAPI string = "https://sparta-gw.battlelog.com/jsonrpc/ps4/api" //交换和行动包查询
)

// ea账密
const (
	USERNAME string = "" //邮箱
	PASSWORD string = ""
)

// 原生API 方法名常量
const (
	//NativeAPI
	ADDVIP       string = "RSP.addServerVip"
	REMOVEVIP    string = "RSP.removeServerVip"
	ADDBAN       string = "RSP.addServerBan"
	REMOVEBAN    string = "RSP.removeServerBan"
	KICK         string = "RSP.kickPlayer"
	STATS        string = "Stats.detailedStatsByPersonaId"
	WEAPONS      string = "Progression.getWeaponsByPersonaId"
	VEHICLES     string = "Progression.getVehiclesByPersonaId"
	PLAYING      string = "GameServer.getServersByPersonaIds"
	RECENTSERVER string = "ServerHistory.mostRecentServers"
	//OperationAPI
	EXCHANGE string = "ScrapExchange.getOffers"
	CAMPAIGN string = "CampaignOperations.getPlayerCampaignStatus"
)

// 游戏代号
const (
	BF1 string = "tunguska"
	BFV string = "casablanca"
	BF4 string = "bf4"
)

// 设置session
var SESSION string = "c3f46e02-ebd3-4b3f-ab5a-e8d548ee14b5"

// bearerAccessToken
var TOKEN string = ""

// Session 获取
func Session(username, password string, refreshToken bool) error {
	if username == "" || password == "" {
		return errors.New("账号信息不完整！")
	}
	login := map[string]interface{}{"username": username, "password": password, "refreshToken": refreshToken}
	//requesting..
	var client = gentleman.New()
	client.URL(SessionAPI)
	client.Use(body.JSON(login))
	//寻找SakuraKooi申请APIKey...
	client.Use(headers.Set("Sakura-Instance-Id", ""))
	client.Use(headers.Set("Sakura-Access-Token", ""))

	res, err := client.Request().Method("POST").Send()
	if err != nil {
		return errors.New("更新session时出错：" + err.Error())
	}
	if gjson.Get(res.String(), "code").Int() != 0 {
		return errors.New("更新session时出错：" + gjson.Get(res.String(), "message").Str)
	}
	var mu sync.Mutex
	mu.Lock()
	SESSION = gjson.Get(res.String(), "data.gatewaySession").Str
	tk := gjson.Get(res.String(), "data.bearerAccessToken").Str
	TOKEN = fmt.Sprintf("%s%s", "Bearer ", tk)
	defer mu.Unlock()
	return nil
}

// NativeAPI 返回json
func ReturnJson(url, method string, parms interface{}) (string, error) {
	var client = gentleman.New()
	client.Use(timeout.Request(time.Second * 30))
	client.URL(url)
	client.Use(body.JSON(parms))
	client.Use(headers.Set("X-Gatewaysession", SESSION))
	res, err := client.Request().Method(method).Send()
	if err != nil {
		return "", errors.New("请求失败")
	}
	data := res.String()
	code := gjson.Get(data, "error.code").Int()
	//如果session过期，重新请求
	if code == -32501 {
		err := Session(USERNAME, PASSWORD, true)
		if err != nil {
			return "", err
		}
		return ReturnJson(url, method, parms)
	}
	return data, nil
}
