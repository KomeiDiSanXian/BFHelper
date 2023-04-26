// Package bf1record 战地相关战绩查询
package bf1record

import (
	"errors"
	"fmt"
	"sort"

	"github.com/tidwall/gjson"

	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
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
	result := gjson.GetMany(data,
		"stats.2.value",
		"stats.4.value",
		"stats.5.value",
		"stats.6.value",
		"stats.7.value",
		"stats.9.value",
		"stats.10.value",
		"stats.11.value",
		"stats.12.value",
		"stats.13.value",
		"stats.14.value",
		"stats.15.value",
		"stats.16.value",
		"stats.17.value",
		"stats.18.value",
		"stats.19.value",
		"stats.20.displayValue",
		"stats.26.value",
		"stats.27.value",
		"stats.31.value",
		"stats.32.value",
		"stats.35.value",
		"stats.37.value",
		"stats.41.value",
		"stats.54.value",
	)
	stat := &Stat{
		SPM:               result[0].Str,
		TotalKD:           result[1].Str,
		WinPercent:        result[2].Str,
		KillsPerGame:      result[3].Str,
		Kills:             result[4].Str,
		Deaths:            result[5].Str,
		KPM:               result[6].Str,
		Losses:            result[7].Str,
		Wins:              result[8].Str,
		InfantryKills:     result[9].Str,
		InfantryKPM:       result[10].Str,
		InfantryKD:        result[11].Str,
		VehicleKills:      result[12].Str,
		VehicleKPM:        result[13].Str,
		Rank:              result[14].Str,
		Skill:             result[15].Str,
		TimePlayed:        result[16].Str,
		MVP:               result[17].Str,
		Accuracy:          result[18].Str,
		DogtagsTaken:      result[19].Str,
		Headshots:         result[20].Str,
		HighestKillStreak: result[21].Str,
		LongestHeadshot:   result[22].Str,
		Revives:           result[23].Str,
		CarriersKills:     result[24].Str,
	}
	return stat, err
}

// GetWeapons 获取武器
func GetWeapons(pid string, class string) (*WeaponSort, error) {
	post := NewPostWeapon(pid)
	data, err := rsp.ReturnJSON(rsp.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	var result []gjson.Result
	if class == ALL {
		result = gjson.Get(data, "result.#.weapons|@flatten").Array()
	} else {
		result = gjson.Get(data, "result.#(categoryId=\""+class+"\").weapons").Array()
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
	post := NewPostVehicle(pid)
	data, err := rsp.ReturnJSON(rsp.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	gets := gjson.Get(data, "result").Array()
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
	post := NewPostStats(pid)
	data, err := rsp.ReturnJSON(rsp.NativeAPI, "POST", post)
	if err != nil {
		return -1, -1, err
	}
	result := gjson.GetMany(data,
		"result.basicStats.kills",
		"result.basicStats.deaths",
		"result.basicStats.kpm",
	)
	kd = result[0].Float() / result[1].Float()
	kpm = result[2].Num
	return kd, kpm, err
}
