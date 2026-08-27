// Package main BFHelper 独立 bot 入口（Remilia 框架）。
//
// 构建：
//
//	go build -o bfhelper.exe .
//
// 配置：编辑 config.yaml 后启动。
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	bfhelper "github.com/KomeiDiSanXian/BFHelper/remilia"
	"github.com/KomeiDiSanXian/remilia/builtin/bundle"
	builtinstorage "github.com/KomeiDiSanXian/remilia/builtin/storage"
	"github.com/KomeiDiSanXian/remilia/config"
	"github.com/KomeiDiSanXian/remilia/core/engine"
	"github.com/KomeiDiSanXian/remilia/infra/storage"
	"github.com/KomeiDiSanXian/remilia/platform"
	"github.com/KomeiDiSanXian/remilia/platform/onebot"
	"github.com/KomeiDiSanXian/remilia/plugin"
	"gopkg.in/yaml.v3"
)

// Config 顶层启动配置（Remilia 框架配置）。
type Config struct {
	Mode    string        `yaml:"mode"`
	Remilia RemiliaConfig `yaml:"remilia"`
}

// RemiliaConfig OneBot 适配器配置。
type RemiliaConfig struct {
	WSServer   string  `yaml:"ws_server"`
	WSToken    string  `yaml:"ws_token"`
	SuperUsers []int64 `yaml:"super_users"`
}

func main() {
	cfg := loadConfig("config.yaml")

	eng := engine.NewEngine()

	botCfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config.yaml: %v", err)
	}

	pm := plugin.NewManager(eng,
		plugin.WithConfigProvider(plugin.NewYAMLConfigProvider(botCfg)),
	)

	if err := pm.Register(builtinstorage.New(
		storage.WithDSN("data/bot.db"),
	)); err != nil {
		log.Fatalf("failed to register storage: %v", err)
	}

	allPlugins := make([]*plugin.Descriptor, 0, len(bundle.Core())+1)
	allPlugins = append(allPlugins, bundle.Core()...)
	allPlugins = append(allPlugins, bfhelper.New())

	regCtx := context.Background()
	if err := pm.RegisterBatch(regCtx, allPlugins); err != nil {
		log.Fatalf("failed to register plugins: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := pm.StartAll(ctx); err != nil {
		cancel()
		log.Fatalf("failed to start plugins: %v", err)
	}
	defer cancel()

	adapterCfg := onebot.DefaultConfig(cfg.Remilia.WSServer)
	adapterCfg.Token = cfg.Remilia.WSToken
	adapter := onebot.NewForwardWSAdapter(adapterCfg)

	go func() {
		if err := adapter.Start(ctx, func(event platform.Event) {
			eng.ProcessPlatformEventEx(event, adapter.Sender(), "")
		}); err != nil {
			log.Printf("onebot adapter stopped: %v", err)
		}
	}()

	log.Println("BFHelper started (Remilia mode)")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	cancel()
	if err := eng.Shutdown(ctx); err != nil {
		log.Printf("engine shutdown: %v", err)
	}
}

func loadConfig(path string) *Config {
	cfg := &Config{
		Mode: "remilia",
		Remilia: RemiliaConfig{
			WSServer: "ws://127.0.0.1:6700",
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
