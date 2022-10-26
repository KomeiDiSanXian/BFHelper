package bf1api

import (
	"errors"
	"fmt"
	"sync"

	"github.com/tidwall/gjson"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
)

// APIs
const (
	AuthAPI      string = "https://api.s-wg.net/ServersCollection/getPlayerAll?DisplayName="
	NativeAPI    string = "https://sparta-gw.battlelog.com/jsonrpc/pc/api"
	SessionAPI   string = "https://battlefield-api.sakurakooi.cyou/account/login"
	OperationAPI string = "https://sparta-gw.battlelog.com/jsonrpc/ps4/api" //交换和行动包查询
	EASBAPI      string = "https://delivery.easb.cc/games/get_server_status"
)

// error code
const (
	ErrServerNotFound string = "-34501"
)

/*
	下面两个必填
*/
// ea账密
const (
	UserName string = "" //邮箱
	Password string = ""
)

// Sakura API info, contacts to SakuraKooi to get it
const (
	SakuraID    string = ""
	SakuraToken string = ""
)

// 原生API 方法名常量
const (
	//NativeAPI
	ADDVIP       string = "RSP.addServerVip"
	REMOVEVIP    string = "RSP.removeServerVip"
	ADDBAN       string = "RSP.addServerBan"
	REMOVEBAN    string = "RSP.removeServerBan"
	KICK         string = "RSP.kickPlayer"
	MAPS         string = "RSP.chooseLeve"
	SERVERDETALS string = "GameServer.getFullServerDetails"
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
var SESSION string = ""

// bearerAccessToken
var TOKEN string = ""

// post operation struct
type post struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Game string `json:"game"`
	} `json:"params"`
	ID string `json:"id"`
}

// unmarshal json
type Pack struct {
	RemainTime int64
	ResetTime  int64
	Name       string
	Desc       string
	Op1Name    string
	Op2Name    string
}

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
	client.Use(headers.Set("Sakura-Instance-Id", SakuraID))
	client.Use(headers.Set("Sakura-Access-Token", SakuraToken))

	res, err := client.Request().Method("POST").Send()
	if err != nil {
		return errors.New("更新session时出错")
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
		err := Session(UserName, Password, true)
		if err != nil {
			return "", err
		}
		return ReturnJson(url, method, parms)
	}
	return data, nil
}

// 查询该周交换
func GetExchange() (map[string][]string, error) {
	post := &post{
		Jsonrpc: "2.0",
		Method:  EXCHANGE,
		Params: struct {
			Game string "json:\"game\""
		}{
			Game: BF1,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
	data, err := ReturnJson(OperationAPI, "POST", post)
	if err != nil {
		return nil, errors.New("获取交换失败")
	}
	var exmap map[string][]string = make(map[string][]string)
	for _, v := range gjson.Get(data, "result.items.#.item").Array() {
		var wpname string = v.Get("parentName").Str
		if wpname == "" {
			wpname = "其他"
		}
		exmap[wpname] = append(exmap[wpname], v.Get("name").Str)
	}
	return exmap, err
}

// 查询本周行动包
func GetCampaignPacks() (*Pack, error) {
	post := &post{
		Jsonrpc: "2.0",
		Method:  CAMPAIGN,
		Params: struct {
			Game string "json:\"game\""
		}{
			Game: BF1,
		},
		ID: "ed26fa43-816d-4f7b-a9d8-de9785ae1bb6",
	}
	data, err := ReturnJson(OperationAPI, "POST", post)
	if err != nil {
		return nil, errors.New("获取行动包失败")
	}
	result := gjson.GetMany(data,
		"result.minutesRemaining",
		"result.name",
		"result.shortDesc",
		"result.op1.name",
		"result.op2.name",
		"result.minutesToDailyReset",
	)
	return &Pack{
		RemainTime: result[0].Int(),
		Name:       result[1].Str,
		Desc:       result[2].Str,
		Op1Name:    result[3].Str,
		Op2Name:    result[4].Str,
		ResetTime:  result[5].Int(),
	}, err
}

// 获取玩家pid
func GetPersonalID(name string) (string, error) {
	cli := gentleman.New()
	cli.URL("https://gateway.ea.com/proxy/identity/personas?namespaceName=cem_ea_id&displayName=" + name)
	cli.Use(headers.Set("X-Expand-Results", "true"))
	cli.Use(headers.Set("Authorization", TOKEN))
	cli.Use(headers.Set("Host", "gateway.ea.com"))
	res, err := cli.Request().Send()
	if err != nil {
		return "", errors.New("获取玩家pid失败")
	}
	info := gjson.Get(res.String(), "error").Str
	if info == "invalid_access_token" || info == "invalid_oauth_info" {
		err := Session(UserName, Password, true)
		if err != nil {
			return "", err
		}
		return GetPersonalID(name)
	}
	if info != "" {
		return "", errors.New(info)
	}
	if gjson.Get(res.String(), "personas.persona.0.personaId").String() == "" {
		return "", errors.New("获取玩家pid失败")
	}
	return gjson.Get(res.String(), "personas.persona.0.personaId").String(), err
}
