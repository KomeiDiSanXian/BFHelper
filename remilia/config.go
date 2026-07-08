package remilia

// Config 是 BFHelper 插件的配置结构体。
//
// 通过 remilia 的 plugin.Config 加载，对应 yaml 中的 bfhelper 配置段。
type Config struct {
	Account AccountConfig `yaml:"account"`
	BFEAC   BFEACConfig   `yaml:"bfeac"`
	Blaze   BlazeConfig   `yaml:"blaze"`
}

// AccountConfig EA 账号认证配置。
type AccountConfig struct {
	// SID EA 登录后的 sid cookie 值。
	SID string `yaml:"sid"`
	// Remid EA 登录后的 remid cookie 值。
	Remid string `yaml:"remid"`
	// Game 目标游戏，可选值: bf1, bfv, bf4。
	Game string `yaml:"game"`
}

// BFEACConfig BFEAC 反作弊查询配置。
type BFEACConfig struct {
	// APIKey BFEAC API 密钥（可选）。
	APIKey string `yaml:"api_key"`
}

// BlazeConfig Blaze TCP 实时推送配置。
type BlazeConfig struct {
	// Enabled 是否全局启用 Blaze TCP 连接。
	Enabled bool `yaml:"enabled"`
}

// DefaultConfig 返回默认配置。
func DefaultConfig() Config {
	return Config{
		Account: AccountConfig{
			Game: "bf1",
		},
		Blaze: BlazeConfig{
			Enabled: false,
		},
	}
}
