package public

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"uuid"

	"github.com/tidwall/gjson"
)

const (
	// SpartaAPIURL BF1/BF4 Sparta JSON-RPC 网关地址。
	SpartaAPIURL = "https://sparta-gw.battlelog.com/jsonrpc/pc/api"
	// spartaOperationURL 交换/行动包查询使用的网关地址（沿用旧实现的 PS4 端点）。
	spartaOperationURL = "https://sparta-gw.battlelog.com/jsonrpc/ps4/api"
	// spartaGame 当前查询的游戏代号。
	spartaGame = "tunguska"
)

// spartaRPC 发送 Sparta JSON-RPC 请求，返回完整响应体。
func spartaRPC(apiURL, session, method string, params map[string]any) (gjson.Result, error) {
	body := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      uuid.New().String(),
	}
	data, err := json.Marshal(body)
	if err != nil {
		return gjson.Result{}, err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(data))
	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gatewaysession", session)
	req.Header.Set("X-HostingGameId", spartaGame)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	res := gjson.ParseBytes(respBody)
	if res.Get("error").Exists() {
		return res, fmt.Errorf("sparta RPC 错误 [%d]: %s",
			res.Get("error.code").Int(), res.Get("error.message").Str)
	}
	return res, nil
}

// Weapon 武器数据。
type Weapon struct {
	Name       string
	Kills      float64
	Accuracy   string
	KPM        string
	Headshots  string
	Efficiency string
}

// GetWeapons 获取玩家武器数据。
//
// category 为武器类别，传入 core.WeaponALL 查询全部武器。
func GetWeapons(pid, apiURL, session, category string) ([]Weapon, error) {
	data, err := spartaRPC(apiURL, session, "Progression.getWeaponsByPersonaId", map[string]any{
		"game":      spartaGame,
		"personaId": pid,
	})
	if err != nil {
		return nil, err
	}

	var result []gjson.Result
	if category == "ALL" {
		result = data.Get("result.#.weapons|@flatten").Array()
	} else {
		result = data.Get(`result.#(categoryId="` + category + `").weapons`).Array()
	}

	ws := make([]Weapon, 0, len(result))
	for i := range result {
		gets := gjson.GetMany(result[i].Raw,
			"name",
			"stats.values.kills",
			"stats.values.headshots",
			"stats.values.accuracy",
			"stats.values.seconds",
			"stats.values.hits",
		)
		kills := gets[1].Float()
		seconds := gets[4].Float()
		heads := gets[2].Float()
		hits := gets[5].Float()

		var kpm, headshots, eff float64
		if kills != 0 && seconds != 0 {
			headshots = heads / kills * 100
			kpm = kills / seconds * 60
			eff = hits / kills
		}
		ws = append(ws, Weapon{
			Name:       gets[0].Str,
			Kills:      kills,
			Accuracy:   fmt.Sprintf("%.2f%%", gets[3].Float()),
			KPM:        fmt.Sprintf("%.3f", kpm),
			Headshots:  fmt.Sprintf("%.2f%%", headshots),
			Efficiency: fmt.Sprintf("%.3f", eff),
		})
	}
	sort.Slice(ws, func(i, j int) bool { return ws[i].Kills > ws[j].Kills })
	return ws, nil
}

// FormatWeapons 格式化武器数据文本。
func FormatWeapons(ws []Weapon, name string) string {
	var sb strings.Builder
	sb.WriteString("id：" + name + "\n")
	for i := 0; i < len(ws) && i < 5; i++ {
		fmt.Fprintf(&sb, "---------------\n武器名：%s\n击杀数：%s\n准度：%s\n爆头率：%s\nKPM：%s\n效率：%s\n",
			ws[i].Name,
			strconv.FormatFloat(ws[i].Kills, 'f', 0, 64),
			ws[i].Accuracy,
			ws[i].Headshots,
			ws[i].KPM,
			ws[i].Efficiency,
		)
	}
	return sb.String()
}

// Vehicle 载具数据。
type Vehicle struct {
	Name      string
	Kills     float64
	Destroyed float64
	KPM       string
	Time      string
}

// GetVehicles 获取玩家载具数据。
func GetVehicles(pid, apiURL, session string) ([]Vehicle, error) {
	data, err := spartaRPC(apiURL, session, "Progression.getVehiclesByPersonaId", map[string]any{
		"game":      spartaGame,
		"personaId": pid,
	})
	if err != nil {
		return nil, err
	}

	gets := data.Get("result").Array()
	vs := make([]Vehicle, 0, len(gets))
	for i := range gets {
		res := gjson.GetMany(gets[i].Raw,
			"name",
			"stats.values.kills",
			"stats.values.destroyed",
			"stats.values.seconds",
		)
		kills := res[1].Float()
		seconds := res[3].Float()
		var kpm float64
		if seconds != 0 {
			kpm = kills / seconds * 60
		}
		vs = append(vs, Vehicle{
			Name:      res[0].Str,
			Kills:     kills,
			Destroyed: res[2].Float(),
			KPM:       fmt.Sprintf("%.3f", kpm),
			Time:      fmt.Sprintf("%.2f", seconds/3600),
		})
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].Kills > vs[j].Kills })
	return vs, nil
}

// FormatVehicles 格式化载具数据文本。
func FormatVehicles(vs []Vehicle, name string) string {
	var sb strings.Builder
	sb.WriteString("id：" + name + "\n")
	for i := range vs {
		fmt.Fprintf(&sb, "------------\n%s\n击杀数：%.0f\tkpm：%s\n击毁数：%.0f\t游玩时间：%s 小时\n",
			vs[i].Name, vs[i].Kills, vs[i].KPM, vs[i].Destroyed, vs[i].Time)
	}
	return sb.String()
}

// GetExchange 查询本期交换皮肤，返回 分类 → 皮肤名列表。
func GetExchange(session string) (map[string][]string, error) {
	data, err := spartaRPC(spartaOperationURL, session, "ScrapExchange.getOffers", map[string]any{
		"game":   spartaGame,
		"offset": 0,
		"limit":  100,
	})
	if err != nil {
		return nil, err
	}

	exchange := make(map[string][]string)
	for _, v := range data.Get("result.items.#.item").Array() {
		parent := v.Get("parentName").Str
		if parent == "" {
			parent = "其他"
		}
		exchange[parent] = append(exchange[parent], v.Get("name").Str)
	}
	return exchange, nil
}

// CampaignPack 行动包信息。
type CampaignPack struct {
	RemainTime int64
	ResetTime  int64
	Name       string
	Desc       string
	Op1Name    string
	Op2Name    string
}

// GetCampaignPacks 查询本期行动包。
func GetCampaignPacks(session string) (*CampaignPack, error) {
	data, err := spartaRPC(spartaOperationURL, session, "CampaignOperations.getPlayerCampaignStatus", map[string]any{
		"game": spartaGame,
	})
	if err != nil {
		return nil, err
	}
	return &CampaignPack{
		RemainTime: data.Get("result.minutesRemaining").Int(),
		Name:       data.Get("result.name").Str,
		Desc:       data.Get("result.shortDesc").Str,
		Op1Name:    data.Get("result.op1.name").Str,
		Op2Name:    data.Get("result.op2.name").Str,
		ResetTime:  data.Get("result.minutesToDailyReset").Int(),
	}, nil
}
