package remilia

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Dev4BF/GoBattlefieldAPI/auth"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/battlefield"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/core"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/public"
	"github.com/KomeiDiSanXian/remilia/infra/kv"
	"github.com/KomeiDiSanXian/remilia/infra/storage"
	"github.com/KomeiDiSanXian/remilia/plugin"
)

// Plugin 是 BFHelper 插件的运行时实例。
type Plugin struct {
	*Config
	core.PlayerStore
	core.GroupStore
	core.ServerStore
	core.AdminStore
	core.BlazeWatchStore
	core.OperationLogStore
	core.DBProvider
	*public.BTRClient
	*public.Bili22Client
	*battlefield.Client
	kv *kv.DB
}

// New 创建 BFHelper 插件描述符。
func New() *plugin.Descriptor {
	p := &Plugin{
		BTRClient:    public.NewBTRClient(),
		Bili22Client: public.NewBili22Client(),
		Client:       battlefield.NewClient(),
	}

	return &plugin.Descriptor{
		Name:         "bfhelper",
		Version:      "2.1.0",
		Deps:         []string{"storage"},
		OptionalDeps: []string{"permission"},
		Meta:         newMetadata(),
		Setup:        p.setup,
		Teardown:     p.teardown,
	}
}

func (p *Plugin) setup(ctx *plugin.SetupContext) (any, error) {
	p.Config = readConfig(ctx)

	store := ctx.Service[*storage.Plugin]("storage")
	db := store.DB()
	if err := db.AutoMigrate(&core.Player{}, &core.Group{}, &core.Server{}, &core.Admin{}, &core.BlazeWatch{}, &core.OperationLog{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	coreStore := core.NewStorage(db)
	p.PlayerStore = coreStore
	p.GroupStore = coreStore
	p.ServerStore = coreStore
	p.AdminStore = coreStore
	p.BlazeWatchStore = coreStore
	p.OperationLogStore = coreStore
	p.DBProvider = coreStore

	if !ctx.DryRun {
		kvDB, err := kv.Open("./data/bfhelper/credentials.kv")
		if err != nil {
			return nil, fmt.Errorf("open kv store: %w", err)
		}
		p.kv = kvDB
		p.tryAutoLogin(ctx)
	}

	p.registerCommands(ctx)
	return p, nil
}

func readConfig(ctx *plugin.SetupContext) *Config {
	cfg := DefaultConfig()
	if account := ctx.Config.GetStringMap("account", nil); account != nil {
		if v, ok := account["sid"].(string); ok && v != "" {
			cfg.Account.SID = v
		}
		if v, ok := account["remid"].(string); ok && v != "" {
			cfg.Account.Remid = v
		}
		if v, ok := account["game"].(string); ok && v != "" {
			cfg.Account.Game = v
		}
	}
	if blaze := ctx.Config.GetStringMap("blaze", nil); blaze != nil {
		if v, ok := blaze["enabled"].(bool); ok {
			cfg.Blaze.Enabled = v
		}
	}
	if v := ctx.Config.GetStringSlice("super_users", nil); v != nil {
		for _, s := range v {
			var qq int64
			if _, err := fmt.Sscanf(s, "%d", &qq); err == nil && qq > 0 {
				cfg.SuperUsers = append(cfg.SuperUsers, qq)
			}
		}
	}
	if bfeac := ctx.Config.GetStringMap("bfeac", nil); bfeac != nil {
		if v, ok := bfeac["api_key"].(string); ok {
			cfg.BFEAC.APIKey = v
		}
	}
	return &cfg
}

func (p *Plugin) teardown(_ *plugin.TeardownContext) error {
	if p.kv != nil {
		return p.kv.Close()
	}
	return nil
}

func (p *Plugin) tryAutoLogin(ctx *plugin.SetupContext) {
	sid, remid := p.Account.SID, p.Account.Remid

	if sid == "" || remid == "" {
		data, err := p.kv.Get([]byte("ea_credentials"))
		if err == nil && len(data) > 0 {
			var creds core.Credentials
			if json.Unmarshal(data, &creds) == nil {
				sid, remid = creds.SID, creds.Remid
			}
		}
	}

	if sid == "" || remid == "" {
		ctx.Log.Info("No EA credentials configured, manual login required")
		return
	}

	eaGame := toEAGame(p.Account.Game)
	if err := p.Login(context.Background(), sid, remid, eaGame); err != nil {
		ctx.Log.Warn(fmt.Sprintf("Auto login failed: %v", err))
		return
	}
	ctx.Log.Info("EA auto login successful")

	game, _ := core.ParseGame(p.Account.Game)
	p.persistCredentials(sid, remid, game)
}

func (p *Plugin) persistCredentials(sid, remid string, game core.Game) {
	if p.kv == nil {
		return
	}
	creds := core.Credentials{SID: sid, Remid: remid, Game: game}
	data, _ := json.Marshal(creds)
	_ = p.kv.Set([]byte("ea_credentials"), data)
}

func toEAGame(gameStr string) auth.Game {
	switch gameStr {
	case "bfv":
		return auth.GameBFV
	case "bf4":
		return auth.GameBF4
	case "bf2042":
		return auth.GameBF2042
	case "bf6":
		return auth.GameBF6
	default:
		return auth.GameBF1
	}
}
