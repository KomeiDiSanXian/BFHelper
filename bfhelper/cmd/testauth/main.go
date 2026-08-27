// Package main 提供 EA 认证与 Sparta 网关的独立调试工具。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"uuid"
)

const (
	EAAuthURL      = "https://accounts.ea.com/connect/auth"
	EARedirectURI  = "nucleus:rest"
	EAClientID     = "ORIGIN_JS_SDK"
	SpartaClientID = "sparta-backend-as-user-pc"
	NativeAPI      = "https://sparta-gw.battlelog.com/jsonrpc/pc/api"
)

func main() {
	sid := flag.String("sid", "", "EA accounts.ea.com sid cookie")
	remid := flag.String("remid", "", "EA accounts.ea.com remid cookie")
	flag.Parse()

	if *sid == "" || *remid == "" {
		*sid = os.Getenv("EA_SID")
		*remid = os.Getenv("EA_REMID")
	}

	if *sid == "" || *remid == "" {
		fmt.Println("错误: 请提供 sid 和 remid")
		fmt.Println("用法: go run main.go -sid=xxx -remid=xxx")
		fmt.Println("或者: set EA_SID=xxx && set EA_REMID=xxx && go run main.go")
		fmt.Println("\n获取方法:")
		fmt.Println("  1. 在浏览器打开 https://www.ea.com/ 并登录")
		fmt.Println("  2. F12 → Application → Cookies → accounts.ea.com")
		fmt.Println("  3. 复制 sid 和 remid 的值")
		os.Exit(1)
	}

	fmt.Println("=== EA OAuth 认证测试 ===")

	// Step 1: 获取 access_token，同时捕获 Set-Cookie 刷新 sid/remid
	fmt.Println("[Step 1/3] 获取 EA access_token...")
	accessToken, resp, err := getEAAccessToken(*sid, *remid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "失败: %v\n", err)
		os.Exit(1)
	}
	// 捕获 Set-Cookie 中刷新的 sid/remid
	for _, c := range resp.Cookies() {
		switch c.Name {
		case "sid":
			*sid = c.Value
		case "remid":
			*remid = c.Value
		}
	}
	_ = resp.Body.Close()
	fmt.Printf("  access_token: %s...%s\n", accessToken[:20], accessToken[len(accessToken)-10:])
	fmt.Printf("  刷新后的 sid: %s...%s\n", (*sid)[:10], (*sid)[len(*sid)-5:])
	fmt.Printf("  刷新后的 remid: %s...%s\n\n", (*remid)[:10], (*remid)[len(*remid)-5:])

	// Step 2: 换取 Sparta auth code
	fmt.Println("[Step 2/3] 获取 Sparta auth code...")
	authCode, err := getSpartaAuthCode(accessToken, *sid, *remid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  auth_code: %s...%s\n\n", authCode[:20], authCode[len(authCode)-10:])

	// Step 3: 换取 gateway session
	fmt.Println("[Step 3/3] 获取 gateway session...")
	sessionID, err := getGatewaySession(authCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  X-Gatewaysession: %s\n\n", sessionID)

	fmt.Println("=== 认证成功! ===")
	fmt.Printf("X-Gatewaysession: %s\n", sessionID)
	fmt.Printf("Token:           Bearer %s...%s\n", accessToken[:20], accessToken[len(accessToken)-10:])

	// 用 session 调一次 RSP.getServerDetails（需 gameId），验证 session 有效
	// 如果返回 -32602 说明是参数错误而非 session 问题；-32501 说明 session 无效
	fmt.Println("\n[验证] 发送测试请求验证 session...")
	err = testSession(sessionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "测试请求失败: %v\n", err)
		fmt.Println("(如果错误码不是 -32501 则 session 依然有效)")
	} else {
		fmt.Println("session 验证通过!")
	}

	fmt.Println("\n将以下配置填入 settings.yml 的 Account 段即可:")
	fmt.Println("  SID:   " + *sid)
	fmt.Println("  Remid: " + *remid)
	fmt.Println("  (Session 和 Token 会在运行时自动获取)")
}

func getEAAccessToken(sid, remid string) (string, *http.Response, error) {
	u := fmt.Sprintf("%s?client_id=%s&response_type=token&redirect_uri=%s&prompt=none&release_type=prod",
		EAAuthURL, EAClientID, url.QueryEscape(EARedirectURI))

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Cookie", fmt.Sprintf("remid=%s; sid=%s", remid, sid))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorCode   string `json:"error_code"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, fmt.Errorf("解析响应失败: %s", string(body))
	}

	if result.Error != "" {
		return "", nil, fmt.Errorf("EA 返回错误: %s (code: %s)", result.Error, result.ErrorCode)
	}
	if result.AccessToken == "" {
		return "", nil, fmt.Errorf("access_token 为空，sid/remid 可能已过期。响应: %s", string(body))
	}
	return result.AccessToken, resp, nil
}

func getSpartaAuthCode(accessToken, sid, remid string) (string, error) {
	u := fmt.Sprintf("%s?access_token=%s&client_id=%s&response_type=code&release_type=prod",
		EAAuthURL, accessToken, SpartaClientID)

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", fmt.Sprintf("remid=%s; sid=%s", remid, sid))
	req.Header.Set("User-Agent", "Mozilla/5.0 EA Download Manager Origin/10.5.94.46774")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("获取 auth code 失败: HTTP %d, body: %s", resp.StatusCode, string(body))
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("重定向响应中缺少 Location 头")
	}

	parsedURL, err := url.Parse(location)
	if err != nil {
		return "", fmt.Errorf("解析重定向 URL 失败: %w", err)
	}
	code := parsedURL.Query().Get("code")
	if code == "" {
		return "", fmt.Errorf("重定向 URL 中缺少 code 参数: %s", location)
	}
	return code, nil
}

func getGatewaySession(authCode string) (string, error) {
	body := map[string]any{
		"jsonrpc": "2.0",
		"method":  "Authentication.getEnvIdViaAuthCode",
		"params": map[string]any{
			"authCode": authCode,
			"locale":   "zh-tw",
		},
		"id": uuid.New().String(),
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, NativeAPI, strings.NewReader(string(bodyJSON)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Guest", "no-session-id")
	req.Header.Set("X-HostingGameId", "tunguska")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Error *struct {
			Code int `json:"code"`
		} `json:"error"`
		Result *struct {
			SessionID string `json:"sessionId"`
		} `json:"result"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %s", string(respBody))
	}

	if result.Error != nil && result.Error.Code != 0 {
		return "", fmt.Errorf("JSON-RPC 错误 code %d, 响应: %s", result.Error.Code, string(respBody))
	}
	if result.Result == nil || result.Result.SessionID == "" {
		return "", fmt.Errorf("sessionId 为空, 响应: %s", string(respBody))
	}
	return result.Result.SessionID, nil
}

func testSession(sessionID string) error {
	body := `{"jsonrpc":"2.0","method":"GameServer.getServerDetails","params":{"game":"tunguska"},"id":"` + uuid.New().String() + `"}`

	req, err := http.NewRequest(http.MethodPost, NativeAPI, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gatewaysession", sessionID)
	req.Header.Set("X-HostingGameId", "tunguska")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		Error *struct {
			Code int `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %s", string(respBody))
	}

	if result.Error != nil && result.Error.Code == -32501 {
		return fmt.Errorf("session 无效 (code -32501)")
	}
	if result.Error != nil && result.Error.Code != 0 {
		return fmt.Errorf("服务器返回错误 code %d", result.Error.Code)
	}
	return nil
}
