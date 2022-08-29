package bf1record

// 武器种类
const (
	ALL     string = "ALL"
	Elite   string = "ID_P_CAT_FIELDKIT"         //精英兵
	LMG     string = "ID_P_CAT_LMG"              //轻机枪
	Melee   string = "ID_P_CAT_MELEE"            //近战武器
	Gadget  string = "ID_P_CAT_GADGET"           //配备
	Semi    string = "ID_P_CAT_SEMI"             //半自动
	Grenade string = "ID_P_CAT_GRENADE"          //手榴弹
	SIR     string = "ID_P_CAT_SIR"              //制式步枪
	Shotgun string = "ID_P_CAT_SHOTGUN"          //霰弹枪
	Dirver  string = "ID_P_CAT_VEHICLEKITWEAPON" //驾驶员
	SMG     string = "ID_P_CAT_SMG"              //冲锋枪
	Sidearm string = "ID_P_CAT_SIDEARM"          //手枪
	Bolt    string = "ID_P_CAT_BOLT"             //步枪
)

type WeaponSort []Weapons
type VehicleSort []Vehicle

func (a WeaponSort) Len() int            { return len(a) }
func (a WeaponSort) Swap(i, j int)       { a[i], a[j] = a[j], a[i] }
func (a WeaponSort) Less(i, j int) bool  { return a[i].Kills > a[j].Kills }
func (a VehicleSort) Len() int           { return len(a) }
func (a VehicleSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VehicleSort) Less(i, j int) bool { return a[i].Kills > a[j].Kills }

//战绩
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
	Skill             string
	TimePlayed        string
	MVP               string
	Accuracy          string
	DogtagsTaken      string
	Headshots         string
	HighestKillStreak string
	LongestHeadshot   string
	Revives           string
	CarriersKills     string
}

//武器
type Weapons struct {
	Name       string
	Kills      float64
	Accuracy   string
	KPM        string
	Headshots  string
	Efficiency string
}

//最近战绩
type Recent []struct {
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

//载具
type Vehicle struct {
	Name      string
	Kills     float64
	Destroyed float64
	KPM       string
	Time      string
}
