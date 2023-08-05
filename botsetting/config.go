// Package botsetting 机器人设置相关
package botsetting

import (
	_ "embed"
	"os"
	"time"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/setting"
	"github.com/sirupsen/logrus"
)

// defaultConfig 默认bot 配置
//
//go:embed default.yaml
var defaultConfig string

// Config 机器人设置
type Config struct {
	BotNames      []string
	CommandPrefix string
	SuperUsers    []int64
	RingLen       uint
	Latency       time.Duration
	MarkMessage   bool
	WSClient      string
	AccessToken   string
}

// Read 读取配置文件
func Read() *Config {
	s, err := setting.NewSetting("botconfig", ".")
	c := Config{}
	if err != nil {
		logrus.Warnln(err)
		generateConfig()
	}
	_ = s.ReadSection("Bot", &c)
	return &c
}

func generateConfig() {
	logrus.Warnln("未找到配置或者出现错误, 正在重新生成...")
	_ = os.WriteFile("botconfig.yaml", []byte(defaultConfig), 0o644)
	logrus.Warnln("配置已生成! 请修改 botconfig.yaml 后重新启动")
	time.Sleep(15 * time.Second)
	os.Exit(0)
}
