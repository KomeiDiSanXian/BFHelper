package bf1record

import (
	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/tidwall/gjson"
)

func GetStats(name string) (*Stat, error) {
	data, err := rsp.ReturnJson("https://battlefieldtracker.com/api/appStats?platform=3&name="+name, "GET", nil)
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
		"stats.32.value",
		"stats.35.value",
		"stats.37.value",
		"stats.41.value",
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
		Headshots:         result[19].Str,
		HighestKillStreak: result[20].Str,
		LongestHeadshot:   result[21].Str,
		Revives:           result[22].Str,
	}
	return stat, err
}
