// Package anticheat BFBan相关
package anticheat

// HackBFBanResp 返回案件信息
type HackBFBanResp struct {
	IsCheater bool
	URL       string
	Status    string
}

// BFBanHackerStatus BFBan 举报状态
var BFBanHackerStatus = map[int]string{
	0: "正在处理",
	1: "实锤",
	2: "等待自证",
	3: "MOSS自证",
	4: "无效举报",
	5: "讨论中",
	6: "需要更多管理员投票",
	8: "刷枪",
}
