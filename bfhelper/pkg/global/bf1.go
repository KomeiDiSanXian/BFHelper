// Package global bf1全局变量
package global

import "github.com/KomeiDiSanXian/BFSession/session"

const (
	EA2788API = "https://ea-api.2788.pro/player?%v=%v" // EA2788API EA2788API地址

	Exchange string = "ScrapExchange.getOffers"                    // Exchange 交换信息获取
	Campaign string = "CampaignOperations.getPlayerCampaignStatus" // Campaign 行动包查询
)

var Session session.Session
