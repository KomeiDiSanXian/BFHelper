package anticheat

type HackEACResp struct {
	URL    string
	Status string
}

// EACHackerStatus BFEAC 举报状态
var EACHackerStatus = map[int]string{
	0: "正在处理",
	1: "实锤",
	2: "证据不足",
	3: "自证通过",
	4: "自证中",
	5: "刷枪",
}
