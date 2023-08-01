// Package player 玩家信息查询
package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/anticheat"
	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq"
	bf1reqbody "github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq/bf1"
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
	data, err := netreq.Request{URL: "https://battlefieldtracker.com/api/appStats?platform=3&name=" + id}.GetRespBodyJSON()
	if err != nil {
		return nil, err
	}
	if !data.IsObject() {
		return nil, errors.New("invalid id")
	}
	stat := &Stat{
		SPM:               data.Get("stats.2.value").Str,
		TotalKD:           data.Get("stats.4.value").Str,
		WinPercent:        data.Get("stats.5.value").Str,
		KillsPerGame:      data.Get("stats.6.value").Str,
		Kills:             data.Get("stats.7.value").Str,
		Deaths:            data.Get("stats.9.value").Str,
		KPM:               data.Get("stats.10.value").Str,
		Losses:            data.Get("stats.11.value").Str,
		Wins:              data.Get("stats.12.value").Str,
		InfantryKills:     data.Get("stats.13.value").Str,
		InfantryKPM:       data.Get("stats.14.value").Str,
		InfantryKD:        data.Get("stats.15.value").Str,
		VehicleKills:      data.Get("stats.16.value").Str,
		VehicleKPM:        data.Get("stats.17.value").Str,
		Rank:              data.Get("stats.18.value").Str,
		Skill:             data.Get("stats.19.value").Str,
		TimePlayed:        data.Get("stats.20.displayValue").Str,
		MVP:               data.Get("stats.26.value").Str,
		Accuracy:          data.Get("stats.27.value").Str,
		DogtagsTaken:      data.Get("stats.31.value").Str,
		Headshots:         data.Get("stats.32.value").Str,
		HighestKillStreak: data.Get("stats.35.value").Str,
		LongestHeadshot:   data.Get("stats.37.value").Str,
		Revives:           data.Get("stats.41.value").Str,
		CarriersKills:     data.Get("stats.54.value").Str,
	}
	return stat, err
}

// GetWeapons 获取武器
func GetWeapons(pid, class string) (*WeaponSort, error) {
	post := bf1reqbody.NewPostWeapon(pid)
	data, err := rsp.ReturnJSON(global.NativeAPI, "POST", post)
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
	post := bf1reqbody.NewPostVehicle(pid)
	data, err := rsp.ReturnJSON(global.NativeAPI, "POST", post)
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

// Get2k battlelog 获取kd,kpm
func Get2k(pid string) (kd float64, kpm float64, err error) {
	post := bf1reqbody.NewPostStats(pid)
	data, err := rsp.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return -1, -1, err
	}
	death := data.Get("result.basicStats.deaths").Float()
	if death == 0 {
		kd = data.Get("result.basicStats.kills").Float()
		return kd, data.Get("result.basicStats.kpm").Float(), nil
	}
	kd = data.Get("result.basicStats.kills").Float() / death
	kpm = data.Get("result.basicStats.kpm").Float()
	return kd, kpm, err
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
