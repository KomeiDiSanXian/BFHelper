// Package bf1api 战地相关api库
package bf1api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// APIs
const (
	AuthAPI      string = "https://api.s-wg.net/ServersCollection/getPlayerAll?DisplayName="
	NativeAPI    string = "https://sparta-gw.battlelog.com/jsonrpc/pc/api"
	SessionAPI   string = "https://battlefield-api.sakurakooi.cyou/account/login"
	OperationAPI string = "https://sparta-gw.battlelog.com/jsonrpc/ps4/api" // 交换和行动包查询
)

// error code
const (
	ErrServerNotFound int64 = -34501
	ErrInvalidMapID   int64 = -32603
	ErrServerOutdate  int64 = -32851
	ErrPlayerIsAdmin  int64 = -32857
	ErrinvalidPlayer  int64 = -32856
	ErrServerNotStart int64 = -32858
)

/*
	下面两个必填
*/
// ea账密
const (
	UserName string = "" // 邮箱
	Password string = ""
)

// Sakura API info, contacts to SakuraKooi to get it
const (
	SakuraID    string = ""
	SakuraToken string = ""
)

// 原生API 方法名常量
const (
	// NativeAPI
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
	SERVERINFO   string = "GameServer.getServerDetails"
	SERVERRSP    string = "RSP.getServerDetails"
	// OperationAPI
	EXCHANGE string = "ScrapExchange.getOffers"
	CAMPAIGN string = "CampaignOperations.getPlayerCampaignStatus"
)

// 游戏代号
const (
	BF1 string = "tunguska"
	BFV string = "casablanca"
	BF4 string = "bf4"
)

var (
	mutex   sync.Mutex
	Session string // Session gatewaysession
	Token   string // Token bearerAccessToken
	Sid     string // Sid cookie sid
	Remid   string // Remid cookie rid
)

// post operation struct
type post struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Game string `json:"game"`
	} `json:"params"`
	ID string `json:"id"`
}

func newpost(method string) *post {
	return &post{
		Jsonrpc: "2.0",
		Method:  method,
		Params: struct {
			Game string "json:\"game\""
		}{
			Game: BF1,
		},
		ID: uuid.NewUUID(),
	}
}

// Pack unmarshal json
type Pack struct {
	RemainTime int64
	ResetTime  int64
	Name       string
	Desc       string
	Op1Name    string
	Op2Name    string
}

// Login 获取 Session token cookies
func Login(username, password string, refreshToken bool) error {
	if username == "" || password == "" {
		return errors.New("账号信息不完整！")
	}
	user := map[string]interface{}{"username": username, "password": password, "refreshToken": refreshToken}
	bodyJSON, err := toJSON(user)
	if err != nil {
		return errors.New("更新session时出错: json marshal error")
	}
	// requesting..
	cli := http.DefaultClient
	req, err := http.NewRequest("POST", SessionAPI, bodyJSON)
	if err != nil {
		return errors.New("更新session时出错: New request failed")
	}
	// 寻找SakuraKooi申请APIKey...
	req.Header.Add("Sakura-Instance-Id", SakuraID)
	req.Header.Add("Sakura-Access-Token", SakuraToken)

	res, err := cli.Do(req)
	if err != nil {
		return errors.New("更新session时出错")
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("更新session时出错: read body error")
	}

	if gjson.GetBytes(resBody, "code").Int() != 0 {
		return errors.New("更新session时出错：" + gjson.GetBytes(resBody, "message").Str)
	}
	datas := gjson.GetManyBytes(resBody, "data.gatewaySession", "data.bearerAccessToken", "data.sid", "data.remid")
	mutex.Lock()
	Session = datas[0].Str
	Token = fmt.Sprintf("%s%s", "Bearer ", datas[1].Str)
	Sid = datas[2].Str
	Remid = datas[3].Str
	mutex.Unlock()
	return nil
}

// ReturnJSON NativeAPI 返回json
func ReturnJSON(url, method string, body interface{}) (string, error) {
	for i := 0; i < 3; i++ { // 3次重试
		data, err := HTTPTry(url, method, body)
		code := gjson.GetBytes(data, "error.code").Int()
		if code == -32501 {
			if err := Login(UserName, Password, true); err != nil {
				logrus.Errorln("[battlefield]", err)
				return "", err
			}
			continue
		}
		if err == nil {
			return string(data), Exception(code)
		}
	}
	return "", errors.New("请求超时，可能是session更新失败")
}

// GetExchange 查询该周交换
func GetExchange() (map[string][]string, error) {
	post := newpost(EXCHANGE)
	data, err := ReturnJSON(OperationAPI, "POST", post)
	if err != nil {
		return nil, errors.New("获取交换失败")
	}
	var exmap = make(map[string][]string)
	for _, v := range gjson.Get(data, "result.items.#.item").Array() {
		var wpname = v.Get("parentName").Str
		if wpname == "" {
			wpname = "其他"
		}
		exmap[wpname] = append(exmap[wpname], v.Get("name").Str)
	}
	return exmap, err
}

// GetCampaignPacks 查询本周行动包
func GetCampaignPacks() (*Pack, error) {
	post := newpost(CAMPAIGN)
	data, err := ReturnJSON(OperationAPI, "POST", post)
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

// GetPersonalID 由name获取玩家pid
func GetPersonalID(name string) (string, error) {
	cli := http.DefaultClient
	url := "https://gateway.ea.com/proxy/identity/personas?namespaceName=cem_ea_id&displayName=" + name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New("获取玩家pid失败")
	}
	req.Header.Add("X-Expand-Results", "true")
	req.Header.Add("Authorization", Token)

	res, err := cli.Do(req)
	if err != nil {
		return "", errors.New("获取玩家pid失败")
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("获取玩家pid失败")
	}
	info := gjson.GetBytes(resBody, "error").Str
	if info == "invalid_access_token" || info == "invalid_oauth_info" {
		err := Login(UserName, Password, true)
		if err != nil {
			return "", err
		}
		return GetPersonalID(name)
	}
	if info != "" {
		return "", errors.New(info)
	}
	pid := gjson.GetBytes(resBody, "personas.persona.0.personaId").String()
	if pid == "" {
		return "", errors.New("获取玩家pid失败")
	}
	return pid, err
}

// Exception 错误码转换
func Exception(errcode int64) error {
	switch errcode {
	case ErrServerNotFound:
		return errors.New("找不到服务器，请检查服务器信息是否正确")
	case ErrInvalidMapID:
		return errors.New("无效的地图id/无权限")
	case ErrServerOutdate:
		return errors.New("找不到服务器/服务器过期")
	case ErrPlayerIsAdmin:
		return errors.New("无权限处理服务器管理")
	case ErrinvalidPlayer:
		return errors.New("找不到该玩家")
	case ErrServerNotStart:
		return errors.New("服务器未开启")
	}
	return nil
}

// any to Reader
func toJSON(data any) (io.Reader, error) {
	buf := &bytes.Buffer{}
	switch data := data.(type) {
	case string:
		buf.WriteString(data)
	case []byte:
		buf.Write(data)
	default:
		if err := json.NewEncoder(buf).Encode(data); err != nil {
			return nil, errors.New("JSON encoding error")
		}
	}
	return io.NopCloser(buf), nil
}

// HTTPTry http请求
func HTTPTry(url, method string, body interface{}) ([]byte, error) {
	cli := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, 5*time.Second) // 5秒连接超时
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(5 * time.Second)) // 5秒接收数据超时
				return conn, nil
			},
		},
	}
	// body is json
	bodyjson, err := toJSON(body)
	if err != nil {
		logrus.Errorln("[battlefield]", err)
		return nil, err
	}

	req, err := http.NewRequest(method, url, bodyjson)
	req.Header.Set("X-Gatewaysession", Session)
	if err != nil {
		logrus.Errorln("[battlefield] newreq err: ", err)
		return nil, errors.New("请求失败")
	}
	res, err := cli.Do(req)
	if err != nil {
		logrus.Errorln("[battlefield] resp err: ", err)
		return nil, errors.New("请求失败，可能是超时了")
	}
	defer res.Body.Close()
	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Errorln("[battlefield]", err)
		return nil, errors.New("解析JSON时发生错误")
	}
	return resbody, nil
}
