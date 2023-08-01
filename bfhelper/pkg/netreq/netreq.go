// Package netreq 网络请求相关
package netreq

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Request 请求结构体
type Request struct {
	Method string
	URL    string
	Header map[string]string
	Body   io.Reader
}

func (r Request) client() *http.Client {
	return &http.Client{
		Timeout: time.Minute,
	}
}

func (r Request) do() (*http.Response, error) {
	if r.Method == "" {
		r.Method = http.MethodGet
	}
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Header {
		req.Header.Set(k, v)
	}
	return r.client().Do(req)
}

func (r Request) respBody() (io.ReadCloser, error) {
	resp, err := r.do()
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// GetRespBodyBytes 向给定URL 发送请求，返回响应体Byte
func (r Request) GetRespBodyBytes() ([]byte, error) {
	b, err := r.respBody()
	if err != nil {
		return nil, err
	}
	defer b.Close()
	return io.ReadAll(b)
}

// GetRespBodyJSON 向给定URL 发送请求，响应转换成JSON
func (r Request) GetRespBodyJSON() (*gjson.Result, error) {
	b, err := r.respBody()
	if err != nil {
		return nil, err
	}
	defer b.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, b)
	if err != nil {
		return nil, err
	}
	result := gjson.Parse(sb.String())
	return &result, nil
}
