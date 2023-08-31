package global

const (
	BFBan = "https://api.gametools.network/bfban/" // BFBan gt联ban api
	BFEAC = "https://api.bfeac.com/"               // BFEAC api
)

// BFEACSetting EAC 设置
type BFEACSetting struct {
	Apikey string
}

// EAC 设置
var EAC BFEACSetting
