package public

import (
	"strings"
	"sync"
)

const (
	bfbanAPI = "https://api.gametools.network/bfban/"
	bfeacAPI = "https://api.bfeac.com/"
)

// BFBan 案件状态。
var bfbanStatus = map[int]string{
	0: "正在处理",
	1: "实锤",
	2: "等待自证",
	3: "MOSS自证",
	4: "无效举报",
	5: "讨论中",
	6: "需要更多管理员投票",
	8: "刷枪",
}

// BFEAC 案件状态。
var bfeacStatus = map[int]string{
	0: "正在处理",
	1: "实锤",
	2: "证据不足",
	3: "自证通过",
	4: "自证中",
	5: "刷枪",
}

// EACInfo BFEAC 案件信息。
type EACInfo struct {
	URL    string
	Status string
}

// BFBanInfo BFBan 案件信息。
type BFBanInfo struct {
	IsCheater bool
	URL       string
	Status    string
}

// CheaterInfo 联ban查询结果。
type CheaterInfo struct {
	EAC   EACInfo
	BFBan BFBanInfo
}

// CheckCheater 查询玩家联ban信息（BFEAC + BFBan 并行）。
func CheckCheater(id string) CheaterInfo {
	var info CheaterInfo
	var wg sync.WaitGroup
	wg.Add(2)

	// BFBan
	go func() {
		defer wg.Done()
		res, err := httpGet(bfbanAPI+"checkban/?names="+id, nil)
		if err != nil || res.Get("errors").String() != "" {
			info.BFBan.Status = "查询失败"
			info.BFBan.URL = "无"
			return
		}
		ban := res.Get("names." + strings.ToLower(id))
		info.BFBan.IsCheater = ban.Get("hacker").Bool()
		if ban.Get("originId").Str != "" {
			info.BFBan.Status = bfbanStatus[int(ban.Get("status").Int())]
			info.BFBan.URL = ban.Get("url").Str
		}
	}()

	// BFEAC
	go func() {
		defer wg.Done()
		res, err := httpGet(bfeacAPI+"case/EAID/"+id, nil)
		if err != nil || res.Get("error_code").Int() != 0 {
			info.EAC.Status = "查询失败: " + res.Get("error_msg").Str
			info.EAC.URL = "无"
			return
		}
		caseInfo := res.Get("data.0")
		info.EAC.Status = bfeacStatus[int(caseInfo.Get("current_status").Int())]
		info.EAC.URL = "https://bfeac.com/?#/case/" + caseInfo.Get("case_id").String()
	}()

	wg.Wait()
	return info
}
