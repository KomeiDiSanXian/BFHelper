// Package player 玩家信息查询
package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/Dev4BF/GoBattlefieldAPI/bf1"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/anticheat"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Cheater 作弊玩家结构体
type Cheater struct {
	EAC   anticheat.HackEACResp
	BFBan anticheat.HackBFBanResp
}

// GetStats 获取战绩信息
func GetStats(id string) (*Stat, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ForceAttemptHTTP2 = false
	transport.TLSClientConfig.NextProtos = []string{"http/1.1"}
	res, err := netreq.Request{
		URL: "https://api.tracker.gg/api/v2/bf1/standard/profile/origin/" + id,
		Header: map[string]string{
			"Host":          "api.tracker.gg",
			"User-Agent":    "Tracker Network App / 3.22.9",
			"Accept":        "application/json; charset=utf-8",
			"x-app-version": "3.22.9",
		},
		Transport: transport,
	}.GetRespBodyJSON()
	if res == nil {
		return nil, errors.New("empty response")
	}
	if err != nil {
		return nil, err
	}
	if res.Get("errors").Exists() {
		return nil, errors.New("invalid id")
	}

	data := res.Get("data.segments.0")
	stat := &Stat{
		SPM:               data.Get("stats.scorePerMinute.displayValue").Str,
		TotalKD:           data.Get("stats.kdRatio.displayValue").Str,
		WinPercent:        data.Get("stats.winPercentage.displayValue").Str,
		KillsPerGame:      data.Get("stats.killsPerRound.displayValue").Str,
		Kills:             data.Get("stats.kills.displayValue").Str,
		Deaths:            data.Get("stats.deaths.displayValue").Str,
		KPM:               data.Get("stats.killsPerMinute.displayValue").Str,
		Losses:            data.Get("stats.losses.displayValue").Str,
		Wins:              data.Get("stats.wins.displayValue").Str,
		InfantryKills:     data.Get("stats.infantryKills.displayValue").Str,
		InfantryKPM:       data.Get("stats.infantryKillsPerMinute.displayValue").Str,
		InfantryKD:        data.Get("stats.infantryKdRatio.displayValue").Str,
		VehicleKills:      data.Get("stats.vehicleKills.displayValue").Str,
		VehicleKPM:        data.Get("stats.vehicleKillsPerMinute.displayValue").Str,
		Rank:              data.Get("stats.rank.displayValue").Str,
		TimePlayed:        data.Get("stats.timePlayed.displayValue").Str,
		MVP:               data.Get("stats.mvp.displayValue").Str,
		Accuracy:          data.Get("stats.shotsAccuracy.displayValue").Str,
		DogtagsTaken:      data.Get("stats.dogtagsTaken.displayValue").Str,
		Headshots:         data.Get("stats.headshots.displayValue").Str,
		HighestKillStreak: data.Get("stats.killStreak.displayValue").Str,
		LongestHeadshot:   data.Get("stats.longestHeadshot.displayValue").Str,
		Revives:           data.Get("stats.revive.displayValue").Str,
	}
	return stat, err
}

// GetWeapons 获取武器
func GetWeapons(pid, class string) (*WeaponSort, error) {
	g := bf1.NewGateway(global.Session.GetSessionID())
	data, err := bf1.GetWeaponsByPersonaID(g, pid)
	if err != nil {
		return nil, err
	}
	var result []gjson.Result
	if class == ALL {
		result = data.Get("result.#.weapons|@flatten").Array()
	} else {
		result = data.Get("result.#(categoryId=\"" + class + "\").weapons").Array()
	}
	return SortWeapon(result), err
}

// GetVehicles 获取载具信息
func GetVehicles(pid string) (*VehicleSort, error) {
	g := bf1.NewGateway(global.Session.GetSessionID())
	data, err := bf1.PlayerVehicles(g, pid)
	if err != nil {
		return nil, err
	}
	gets := data.Get("result").Array()
	var vehicle VehicleSort
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
		vehicle = append(vehicle, Vehicle{
			Name:      res[0].Str,
			Kills:     kills,
			Destroyed: res[2].Float(),
			KPM:       fmt.Sprintf("%.3f", kpm),
			Time:      fmt.Sprintf("%.2f", seconds/3600),
		})
	}
	sort.Sort(vehicle)
	return &vehicle, err
}

// IsHacker 借助id 获取举报信息
func IsHacker(id string) *Cheater {
	var c Cheater
	var wg sync.WaitGroup
	wg.Add(2)
	// bfban
	go func() {
		defer wg.Done()
		result, err := netreq.Request{URL: global.BFBan + "checkban/?names=" + id}.GetRespBodyJSON()
		if err != nil || result.Get("errors").String() != "" {
			c.BFBan.Status = "查询失败"
			c.BFBan.URL = "无"
			return
		}
		bfban := result.Get("names." + strings.ToLower(id))
		c.BFBan.IsCheater = bfban.Get("hacker").Bool()
		if bfban.Get("originId").Str != "" {
			c.BFBan.Status = anticheat.BFBanHackerStatus[int(bfban.Get("status").Int())]
			c.BFBan.URL = bfban.Get("url").Str
		}
	}()
	// bfeac
	go func() {
		defer wg.Done()
		result, err := netreq.Request{URL: global.BFEAC + "case/EAID/" + id}.GetRespBodyJSON()
		if err != nil || result.Get("error_code").Int() != 0 {
			c.EAC.Status = "查询失败: " + result.Get("error_msg").Str
			c.EAC.URL = "无"
			return
		}
		bfeac := result.Get("data.0")
		c.EAC.Status = anticheat.EACHackerStatus[int(bfeac.Get("current_status").Int())]
		c.EAC.URL = "https://bfeac.com/?#/case/" + bfeac.Get("case_id").String()
	}()
	wg.Wait()
	return &c
}

// GetBF1Recent 获取bf1最近战绩
func GetBF1Recent(id string) (result *Recent, err error) {
	u := "https://api.bili22.me/bf1/recent?name=" + id
	data, err := netreq.Request{URL: u}.GetRespBodyBytes()
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	if err != nil {
		return nil, errors.New("ERR: JSON decode failed")
	}
	return result, err
}
