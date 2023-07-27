// Package global bf1全局变量
package global

const (
	NativeAPI    string = "https://sparta-gw.battlelog.com/jsonrpc/pc/api"       // NativeAPI EA JSONRPC API
	SessionAPI   string = "https://battlefield-api.sakurakooi.dev/account/login" // SessionAPI 通过SakuraKooi 获取session 信息
	OperationAPI string = "https://sparta-gw.battlelog.com/jsonrpc/ps4/api"      // OperationAPI 交换和行动包查询
)

const (
	AddVIP        string = "RSP.addServerVip"                   // AddVIP 			单服务器添加VIP
	RemoveVIP     string = "RSP.removeServerVip"                // RemoveVIP 		单服务器移除VIP
	AddBan        string = "RSP.addServerBan"                   // AddBan 			单服务器添加玩家进入ban列
	RemoveBan     string = "RSP.removeServerBan"                // RemoveBan 		单服务器ban列移除玩家
	Kick          string = "RSP.kickPlayer"                     // Kick 			单服务器踢出玩家
	ChooseMap     string = "RSP.chooseLeve"                     // ChooseMap 		单服务器切换地图
	ServerDetails string = "GameServer.getFullServerDetails"    // ServerDetails 	单服务器完整信息查询
	Stats         string = "Stats.detailedStatsByPersonaId"     // Stats 			单玩家战绩获取
	Weapons       string = "Progression.getWeaponsByPersonaId"  // Weapons 			单玩家武器获取
	Vehicles      string = "Progression.getVehiclesByPersonaId" // Vehicles 		单玩家载具获取
	Playing       string = "GameServer.getServersByPersonaIds"  // Playing 			多玩家正在游玩获取
	RecentServer  string = "ServerHistory.mostRecentServers"    // RecentServer 	单玩家游玩服务器历史获取
	ServerInfo    string = "GameServer.getServerDetails"        // ServerInfo 		单服务器游戏信息获取
	ServerRSP     string = "RSP.getServerDetails"               // ServerRSP 		单服务器RSP信息获取

	Exchange string = "ScrapExchange.getOffers"                    // Exchange 交换信息获取
	Campaign string = "CampaignOperations.getPlayerCampaignStatus" // Campaign 行动包查询
)
