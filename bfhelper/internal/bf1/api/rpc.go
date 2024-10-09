// Package bf1api 战地相关api库
package bf1api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Dev4BF/GoBattlefieldAPI/bf1"
	"github.com/Dev4BF/GoBattlefieldAPI/dto"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Pack unmarshal json
type Pack struct {
	RemainTime int64
	ResetTime  int64
	Name       string
	Desc       string
	Op1Name    string
	Op2Name    string
}

func param(m string) *dto.GenericRequest {
	return &dto.GenericRequest{
		Jsonrpc: "2.0",
		Method:  m,
		Params:  map[string]string{"game": global.BF1},
	}
}

// GetExchange 查询该周交换
func GetExchange() (map[string][]string, error) {
	g := bf1.NewGateway(global.Session.GetSessionID())
	data, err := g.Post(param(global.Exchange))
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
	g := bf1.NewGateway(global.Session.GetSessionID())
	data, err := g.Post(param(global.Campaign))
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
	resp, err := http.Get(fmt.Sprintf(global.EA2788API, "name", name))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("获取玩家ID失败")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	jsonResult := gjson.ParseBytes(b)
	if !jsonResult.Exists() {
		return "", errors.New("invalid JSON data")
	}

	if pid, ok := jsonResult.Map()["personaId"]; ok {
		return pid.String(), nil
	}
	return "", errors.New("获取玩家ID失败")
}

// GetServerFullInfo 获取服务器完整信息
func GetServerFullInfo(gameID string) (*gjson.Result, error) {
	r, err := bf1.GetServerFullDetails(bf1.NewGateway(global.Session.GetSessionID()), gameID)
	return &r, err
}
