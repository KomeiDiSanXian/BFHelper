// Package public 提供 BFHelper 的公开数据查询能力。
//
// 本包不依赖任何私有库，仅使用标准库与 gjson：
//   - BTR (tracker.gg) 战绩查询
//   - Bili22 最近战绩查询
//   - BFEAC / BFBan 联ban查询
//   - Sparta 网关武器 / 载具 / 交换 / 行动包查询
package public

import (
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// httpGet 发起 GET 请求并解析响应为 JSON。
func httpGet(url string, header map[string]string) (gjson.Result, error) {
	return doRequest(http.MethodGet, url, nil, header)
}

// doRequest 发起 HTTP 请求并解析响应为 JSON。
func doRequest(method, url string, body io.Reader, header map[string]string) (gjson.Result, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return gjson.Result{}, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(data), nil
}
