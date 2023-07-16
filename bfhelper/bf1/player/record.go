// Package player 战地相关战绩查询
package player

import (
	"errors"
	"fmt"
	"sort"

	"github.com/tidwall/gjson"

	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/global"
	bf1reqbody "github.com/KomeiDiSanXian/BFHelper/bfhelper/netreq/bf1"
)

// GetStats 获取战绩信息
func GetStats(name string) (*Stat, error) {
	if name == "" {
		return nil, errors.New("ID cannot be empty")
	}
	data, err := rsp.ReturnJSON("https://battlefieldtracker.com/api/appStats?platform=3&name="+name, "GET", nil)
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
func GetWeapons(pid string, class string) (*WeaponSort, error) {
	post := bf1reqbody.NewPostWeapon(pid)
	data, err := rsp.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	var result []gjson.Result
	if class == ALL {
		result = data.Get( "result.#.weapons|@flatten").Array()
	} else {
		result = data.Get("result.#(categoryId=\""+class+"\").weapons").Array()
	}
	return SortWeapon(result), err
}

// SortWeapon 武器排序
func SortWeapon(weapons []gjson.Result) *WeaponSort {
	wp := WeaponSort{}
	for i := range weapons {
		gets := gjson.GetMany(weapons[i].Raw,
			"name",
			"stats.values.kills",
			"stats.values.headshots",
			"stats.values.accuracy",
			"stats.values.seconds",
			"stats.values.hits",
			"stats.values.shots",
		)
		kills := gets[1].Float()
		seconds := gets[4].Float()
		heads := gets[2].Float()
		hits := gets[5].Float()
		var (
			kpm       float64
			headshots float64
			eff       float64
		)
		// 除以0的情况
		if kills == 0 || seconds == 0 {
			kpm = 0
			headshots = 0
			eff = 0
		} else {
			headshots = heads / kills * 100
			kpm = kills / seconds * 60
			eff = hits / kills
		}
		wp = append(wp, Weapons{
			Name:       gets[0].Str,
			Kills:      kills,
			Accuracy:   fmt.Sprintf("%.2f%%", gets[3].Float()),
			KPM:        fmt.Sprintf("%.3f", kpm),
			Headshots:  fmt.Sprintf("%.2f%%", headshots),
			Efficiency: fmt.Sprintf("%.3f", eff),
		})
	}
	sort.Sort(wp)
	return &wp
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
	kd = data.Get("result.basicStats.kills").Float() / data.Get("result.basicStats.deaths").Float()
	kpm = data.Get("result.basicStats.kpm").Float()
	return kd, kpm, err
}
