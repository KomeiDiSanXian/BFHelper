//go:build zerobot

// Package main BFHelper ZeroBot 模式入口。
//
// 构建：
//
//	go build -tags zerobot -o bfhelper.exe .
//
// 配置：编辑 config.yaml 后启动。
package main

import (
	"log"
	"os"

	_ "github.com/KomeiDiSanXian/BFHelper/zerobot"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"gopkg.in/yaml.v3"
)

type ZeroBotConfig struct {
	BotNames      []string `yaml:"bot_names"`
	CommandPrefix string   `yaml:"command_prefix"`
	SuperUsers    []int64  `yaml:"super_users"`
	WSServer      string   `yaml:"ws_server"`
	WSToken       string   `yaml:"ws_token"`
}

type Config struct {
	Mode    string         `yaml:"mode"`
	ZeroBot ZeroBotConfig  `yaml:"zerobot"`
}

func main() {
	cfg := loadZerobotConfig("config.yaml")

	var c zero.Config
	c.NickName = cfg.ZeroBot.BotNames
	if len(c.NickName) == 0 {
		c.NickName = []string{"蕾米"}
	}
	c.CommandPrefix = cfg.ZeroBot.CommandPrefix
	if c.CommandPrefix == "" {
		c.CommandPrefix = "."
	}
	c.SuperUsers = cfg.ZeroBot.SuperUsers
	c.RingLen = 40

	wsURL := cfg.ZeroBot.WSServer
	if wsURL == "" {
		wsURL = "ws://127.0.0.1:6700"
	}
	c.Driver = []zero.Driver{
		driver.NewWebSocketClient(wsURL, cfg.ZeroBot.WSToken),
	}

	log.Println("BFHelper started (ZeroBot mode)")
	zero.RunAndBlock(&c, func() {})
}

func loadZerobotConfig(path string) *Config {
	cfg := &Config{
		Mode: "zerobot",
		ZeroBot: ZeroBotConfig{
			BotNames:      []string{"蕾米"},
			CommandPrefix: ".",
			WSServer:      "ws://127.0.0.1:6700",
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("config not found at %s, using defaults", path)
		return cfg
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	return cfg
}
