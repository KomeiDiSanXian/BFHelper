// Package anticheat BFEAC相关
package anticheat

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq"
)

// HackEACResp 返回案件信息
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

// Report 举报信息
type Report struct {
	Target   string        `json:"target_EAID"`
	Body     template.HTML `json:"case_body"`
	GameType int           `json:"game_type"`
	UploadBy Reporter      `json:"report_by"`
}

// Reporter 举报人信息
type Reporter struct {
	Platform  string `json:"report_platform"`
	UserID    string `json:"user_id"`
	Timestamp string
	Details   map[string]any
}

// NewBF1Report 生成一个bf1举报信息
func NewBF1Report(susName string, body template.HTML) Report {
	return Report{
		Target:   susName,
		Body:     body,
		GameType: 1,
	}
}

// WithReporter 附加上举报人信息
func (r *Report) WithReporter(reporter Reporter) Report {
	r.UploadBy = reporter
	return *r
}

// UploadReport 用户进行举报
func UploadReport(repo Report) error {
	jsonData, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	_, err = netreq.Request{
		Method: http.MethodPost,
		URL:    global.BFEAC + "inner_api/case_report",
		Header: map[string]string{
			"Content-Type": "application/json",
			"apikey":       global.BFEACSetting.APIKey,
		},
		Body: bytes.NewReader(jsonData),
	}.GetRespBodyJSON()
	if err != nil {
		return err
	}
	return nil
}
