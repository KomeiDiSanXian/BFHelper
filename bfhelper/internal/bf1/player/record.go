// Package player 战地相关战绩查询
package player

import (
	"fmt"
	"sort"

	"github.com/tidwall/gjson"
)

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
