// Package bf1api 战地相关api库
package bf1api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/netreq"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// error code
//
// TODO: refactor this
const (
	ErrServerNotFound int64 = -34501
	ErrInvalidMapID   int64 = -32603
	ErrServerOutdate  int64 = -32851
	ErrPlayerIsAdmin  int64 = -32857
	ErrinvalidPlayer  int64 = -32856
	ErrServerNotStart int64 = -32858
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
			Game: global.BF1,
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
	result, err := netreq.Request{
		Method: http.MethodPost,
		URL:    global.SessionAPI,
		Header: map[string]string{
			"Sakura-Instance-Id":  global.SakuraAPI.SakuraID,
			"Sakura-Access-Token": global.SakuraAPI.SakuraToken,
		},
		Body: bodyJSON,
	}.GetRespBodyJSON()
	global.Account.Info.Session = result.Get("data.gatewaySession").Str
	global.Account.Info.Token = fmt.Sprintf("%s%s", "Bearer ", result.Get("data.bearerAccessToken").Str)
	global.Account.Info.SID = result.Get("data.sid").Str
	global.Account.Info.Remid = result.Get("data.remid").Str
	return nil
}

// ReturnJSON NativeAPI 返回json
func ReturnJSON(url, method string, body interface{}) (*gjson.Result, error) {
	for i := 0; i < 3; i++ { // 3次重试
		// body is json
		bodyjson, err := toJSON(body)
		if err != nil {
			logrus.Errorln("[battlefield]", err)
			return nil, err
		}
		result, err := netreq.Request{
			Method: method,
			URL:    url,
			Header: map[string]string{"X-Gatewaysession": global.Account.Info.Session},
			Body:   bodyjson,
		}.GetRespBodyJSON()

		code := result.Get("error.code").Int()
		if code == -32501 {
			if err := Login(global.Account.LoginedUser.Username, global.Account.LoginedUser.Password, true); err != nil {
				logrus.Errorln("[battlefield]", err)
				return nil, err
			}
			continue
		}
		if err == nil {
			return result, Exception(code)
		}
	}
	return nil, errors.New("请求超时，可能是session更新失败")
}

// GetExchange 查询该周交换
func GetExchange() (map[string][]string, error) {
	post := newpost(global.Exchange)
	data, err := ReturnJSON(global.OperationAPI, http.MethodPost, post)
	if err != nil {
		return nil, errors.New("获取交换失败")
	}
	var exmap = make(map[string][]string)
	for _, v := range data.Get("result.items.#.item").Array() {
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
	post := newpost(global.Campaign)
	data, err := ReturnJSON(global.OperationAPI, http.MethodPost, post)
	if err != nil {
		return nil, errors.New("获取行动包失败")
	}
	return &Pack{
		RemainTime: data.Get("result.minutesRemaining").Int(),
		Name:       data.Get("result.name").Str,
		Desc:       data.Get("result.shortDesc").Str,
		Op1Name:    data.Get("result.op1.name").Str,
		Op2Name:    data.Get("result.op2.name").Str,
		ResetTime:  data.Get("result.minutesToDailyReset").Int(),
	}, err
}

// GetPersonalID 由name获取玩家pid
func GetPersonalID(name string) (string, error) {
	result, err := netreq.Request{
		Method: http.MethodGet,
		URL:    "https://gateway.ea.com/proxy/identity/personas?namespaceName=cem_ea_id&displayName=" + name,
		Header: map[string]string{
			"X-Expand-Results": "true",
			"Authorization":    global.Account.Info.Token,
		},
	}.GetRespBodyJSON()
	info := result.Get("error").Str
	if info == "invalid_access_token" || info == "invalid_oauth_info" {
		err := Login(global.Account.LoginedUser.Username, global.Account.LoginedUser.Password, true)
		if err != nil {
			return "", err
		}
		return GetPersonalID(name)
	}
	if info != "" {
		return "", errors.New(info)
	}
	pid := result.Get("personas.persona.0.personaId").String()
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
