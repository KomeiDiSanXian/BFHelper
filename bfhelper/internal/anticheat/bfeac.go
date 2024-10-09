// Package anticheat BFEAC相关
package anticheat

import (
	// "bytes"
	_ "embed"
	// "encoding/json"
	// "fmt"
	// "html/template"
	// "net/http"
	// "time"

	// "github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	// "github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq"
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

// ReportHTML 举报html模板
//
//go:embed template.html
var ReportHTML string

// Report 举报信息
type Report struct {
	Target   string   `json:"target_EAID"`
	Body     string   `json:"case_body"`
	GameType int      `json:"game_type"`
	UploadBy Reporter `json:"report_by"`
}

// Reporter 举报人信息
type Reporter struct {
	Platform  string `json:"report_platform"`
	UserID    int64  `json:"user_id"`
	Timestamp int64
	Details   map[string]any
}

// ReportBody 用于举报上传的html模板
type ReportBody struct {
	Link     string
	ImageURL string
	Words    []string
	SelfID   int64
}

// // NewBF1Report 生成一个bf1举报信息
// func NewBF1Report(susName string, body ReportBody) (*Report, error) {
// 	tmplObj := template.Must(template.New("report").Parse(ReportHTML))
// 	var buf bytes.Buffer

// 	err := tmplObj.Execute(&buf, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Report{
// 		Target:   susName,
// 		Body:     buf.String(),
// 		GameType: 1,
// 	}, nil
// }

// // NewReporter 生成举报人信息
// func NewReporter(userID int64) *Reporter {
// 	return &Reporter{
// 		Platform:  "qq",
// 		UserID:    userID,
// 		Timestamp: time.Now().Unix(),
// 		Details:   make(map[string]any),
// 	}
// }

// // WithDetails 给举报人添加更多信息
// func (r *Reporter) WithDetails(detail map[string]any) *Reporter {
// 	rr := *r
// 	rr.Details = r.Details
// 	for k, v := range detail {
// 		rr.Details[k] = v
// 	}
// 	return &rr
// }

// // WithReporter 附加上举报人信息
// func (r *Report) WithReporter(reporter *Reporter) Report {
// 	r.UploadBy = *reporter
// 	return *r
// }

// // UploadReport 用户进行举报
// func (r Report) UploadReport() error {
// 	jsonData, err := json.Marshal(r)
// 	if err != nil {
// 		return err
// 	}
// 	data, err := netreq.Request{
// 		Method: http.MethodPost,
// 		URL:    global.BFEAC + "inner_api/case_report",
// 		Header: map[string]string{
// 			"Content-Type": "application/json",
// 			"apikey":       global.BFEACSetting.APIKey,
// 		},
// 		Body: bytes.NewReader(jsonData),
// 	}.GetRespBodyJSON()
// 	if err != nil {
// 		return err
// 	}
// 	// TODO: 等eac 新后端再做解析
// 	fmt.Println(data.Map())
// 	return nil
// }
