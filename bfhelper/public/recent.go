package public

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Bili22Client 查询 Bili22 最近战绩。
type Bili22Client struct{}

// NewBili22Client 创建 Bili22 客户端。
func NewBili22Client() *Bili22Client { return &Bili22Client{} }

// RecentEntry 最近一场对局记录。
type RecentEntry struct {
	Server   string  `json:"server"`
	Map      string  `json:"map"`
	Mode     string  `json:"mode"`
	Date     int64   `json:"date"`
	Score    int     `json:"score"`
	Kill     int     `json:"kill"`
	Death    int     `json:"death"`
	Kd       float64 `json:"kd"`
	Kpm      float64 `json:"kpm"`
	Accuracy float64 `json:"accuracy"`
	Time     int     `json:"time"`
}

// Recent 最近战绩列表。
type Recent []RecentEntry

// GetRecent 获取玩家最近战绩。
func (c *Bili22Client) GetRecent(name string) (Recent, error) {
	res, err := httpGet("https://api.bili22.me/bf1/recent?name="+name, nil)
	if err != nil {
		return nil, err
	}
	var entries Recent
	if err := json.Unmarshal([]byte(res.Raw), &entries); err != nil {
		return nil, fmt.Errorf("解析最近战绩失败: %w", err)
	}
	return entries, nil
}

// FormatRecent 格式化最近战绩文本。
func FormatRecent(entries Recent, name string) string {
	var sb strings.Builder
	sb.WriteString("id：" + name + "\n")
	for i := range entries {
		e := &entries[i]
		server := e.Server
		if len(server) > 24 {
			server = server[:24]
		}
		sb.WriteString("服务器：" + server + "\n")
		sb.WriteString("地图：" + e.Map + "   (" + e.Mode + ")\n")
		sb.WriteString("kd：" + strconv.FormatFloat(e.Kd, 'f', -1, 64) + "\n")
		sb.WriteString("kpm：" + strconv.FormatFloat(e.Kpm, 'f', -1, 64) + "\n")
		sb.WriteString("游玩时长：" + strconv.Itoa(e.Time/60) + "分钟\n")
		sb.WriteString("---------------\n")
	}
	return sb.String()
}
