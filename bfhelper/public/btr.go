package public

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// BTRClient 查询 BTR（tracker.gg）战绩。
type BTRClient struct{}

// NewBTRClient 创建 BTR 客户端。
func NewBTRClient() *BTRClient { return &BTRClient{} }

// Stat 玩家战绩。
type Stat struct {
	SPM               string
	TotalKD           string
	WinPercent        string
	KillsPerGame      string
	Kills             string
	Deaths            string
	KPM               string
	Losses            string
	Wins              string
	InfantryKills     string
	InfantryKPM       string
	InfantryKD        string
	VehicleKills      string
	VehicleKPM        string
	Rank              string
	TimePlayed        string
	MVP               string
	Accuracy          string
	DogtagsTaken      string
	Headshots         string
	HighestKillStreak string
	LongestHeadshot   string
	Revives           string
}

// GetStats 获取玩家 BTR 战绩。
func (c *BTRClient) GetStats(name string) (*Stat, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ForceAttemptHTTP2 = false
	transport.TLSClientConfig.NextProtos = []string{"http/1.1"}

	req, err := http.NewRequest(http.MethodGet,
		"https://api.tracker.gg/api/v2/bf1/standard/profile/origin/"+name, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "api.tracker.gg")
	req.Header.Set("User-Agent", "Tracker Network App / 3.22.9")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("x-app-version", "3.22.9")

	client := &http.Client{Timeout: 60 * time.Second, Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := gjson.ParseBytes(data)
	if res.Get("errors").Exists() {
		return nil, fmt.Errorf("无效的玩家名: %s", name)
	}

	seg := res.Get("data.segments.0")
	return &Stat{
		SPM:               seg.Get("stats.scorePerMinute.displayValue").Str,
		TotalKD:           seg.Get("stats.kdRatio.displayValue").Str,
		WinPercent:        seg.Get("stats.winPercentage.displayValue").Str,
		KillsPerGame:      seg.Get("stats.killsPerRound.displayValue").Str,
		Kills:             seg.Get("stats.kills.displayValue").Str,
		Deaths:            seg.Get("stats.deaths.displayValue").Str,
		KPM:               seg.Get("stats.killsPerMinute.displayValue").Str,
		Losses:            seg.Get("stats.losses.displayValue").Str,
		Wins:              seg.Get("stats.wins.displayValue").Str,
		InfantryKills:     seg.Get("stats.infantryKills.displayValue").Str,
		InfantryKPM:       seg.Get("stats.infantryKillsPerMinute.displayValue").Str,
		InfantryKD:        seg.Get("stats.infantryKdRatio.displayValue").Str,
		VehicleKills:      seg.Get("stats.vehicleKills.displayValue").Str,
		VehicleKPM:        seg.Get("stats.vehicleKillsPerMinute.displayValue").Str,
		Rank:              seg.Get("stats.rank.displayValue").Str,
		TimePlayed:        seg.Get("stats.timePlayed.displayValue").Str,
		MVP:               seg.Get("stats.mvp.displayValue").Str,
		Accuracy:          seg.Get("stats.shotsAccuracy.displayValue").Str,
		DogtagsTaken:      seg.Get("stats.dogtagsTaken.displayValue").Str,
		Headshots:         seg.Get("stats.headshots.displayValue").Str,
		HighestKillStreak: seg.Get("stats.killStreak.displayValue").Str,
		LongestHeadshot:   seg.Get("stats.longestHeadshot.displayValue").Str,
		Revives:           seg.Get("stats.revive.displayValue").Str,
	}, nil
}

// FormatStat 格式化战绩文本。
func FormatStat(stat *Stat, name string) string {
	return "id：" + name + "\n" +
		"等级：" + stat.Rank + "\n" +
		"游玩时长：" + stat.TimePlayed + "\n" +
		"总kd：" + stat.TotalKD + "(" + stat.Kills + "/" + stat.Deaths + ")" + "\n" +
		"总kpm：" + stat.KPM + "\n" +
		"准度：" + stat.Accuracy + "\n" +
		"爆头数：" + stat.Headshots + "\n" +
		"胜率：" + stat.WinPercent + "(" + stat.Wins + "/" + stat.Losses + ")" + "\n" +
		"场均击杀：" + stat.KillsPerGame + "\n" +
		"步战kd：" + stat.InfantryKD + "\n" +
		"步战击杀：" + stat.InfantryKills + "\n" +
		"步战kpm：" + stat.InfantryKPM + "\n" +
		"载具击杀：" + stat.VehicleKills + "\n" +
		"载具kpm：" + stat.VehicleKPM + "\n" +
		"近战击杀：" + stat.DogtagsTaken + "\n" +
		"最高连杀：" + stat.HighestKillStreak + "\n" +
		"最远爆头：" + stat.LongestHeadshot + "\n" +
		"MVP数：" + stat.MVP + "\n" +
		"作为神医拉起了 " + stat.Revives + " 人"
}
