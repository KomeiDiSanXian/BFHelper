package global

const (
	BFBan = "https://api.gametools.network/bfban/" // BFBan gt联ban api
	BFEAC = "https://api.bfeac.com/"               // BFEAC api
)

type BFEACSetting struct {
	Apikey string
}

var EAC BFEACSetting
