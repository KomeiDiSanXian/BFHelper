// Package player 玩家信息查询
package player

import (
	"errors"
	"fmt"
	"sort"

	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/global"
	bf1reqbody "github.com/KomeiDiSanXian/BFHelper/bfhelper/netreq/bf1"
	"github.com/tidwall/gjson"
)

// Info 玩家信息结构体
type Info struct {
	Name       string
	PersonalID string
}

// NewInfoByQID 根据qq 生成Info
func NewInfoByQID(qid int64) *Info {
	// 先在数据库中查询
	p, err := bf1model.NewPlayerRepository(global.DB).GetByQID(qid)
	// 找不到直接返回nil
	if err != nil {
		return nil
	}
	return &Info{
		Name:       p.DisplayName,
		PersonalID: p.PersonalID,
	}
}

// NewInfoByName 根据给定的玩家名生成Player
func NewInfoByName(name string) *Info {
	player := &Info{Name: name}
	p, err := bf1model.NewPlayerRepository(global.DB).GetByName(name)
	if err != nil {
		// 找不到就查pid
		player.PersonalID, err = rsp.GetPersonalID(name)
		if err != nil {
			return nil
		}
		return player
	}
	player.PersonalID = p.PersonalID
	return player
}

// GetStats 获取战绩信息
func (p *Info) GetStats() (*Stat, error) {
	if p == nil {
		return nil, errors.New("找不到查询的玩家")
	}
	data, err := rsp.ReturnJSON("https://battlefieldtracker.com/api/appStats?platform=3&name="+p.Name, "GET", nil)
	if err != nil {
		return nil, err
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
func (p *Info) GetWeapons(class string) (*WeaponSort, error) {
	if p == nil {
		return nil, errors.New("找不到查询的玩家")
	}
	post := bf1reqbody.NewPostWeapon(p.PersonalID)
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
func (p *Info) GetVehicles() (*VehicleSort, error) {
	if p == nil {
		return nil, errors.New("找不到查询的玩家")
	}
	post := bf1reqbody.NewPostVehicle(p.PersonalID)
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
func (p *Info) Get2k() (kd float64, kpm float64, err error) {
	if p == nil {
		return -1, -1, errors.New("找不到查询的玩家")
	}
	post := bf1reqbody.NewPostStats(p.PersonalID)
	data, err := rsp.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return -1, -1, err
	}
	kd = data.Get("result.basicStats.kills").Float() / data.Get("result.basicStats.deaths").Float()
	kpm = data.Get("result.basicStats.kpm").Float()
	return kd, kpm, err
}
