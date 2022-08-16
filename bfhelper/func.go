package bfhelper

import (
	"encoding/json"
	"errors"
	"io/ioutil"	
	"os"
	"strings"

	rsp "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/tidwall/gjson"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
)

//字典
var twmap map[string]string

//初始化
func init() {
	//读字典
	f, err := os.Open(engine.DataFolder() + "dic/dic.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &twmap)
	if err != nil {
		panic(err)
	}
}

//查询是否被实锤为外挂
func IsGetBan(id string) bool {
	cli := gentleman.New()
	cli.URL("https://api.gametools.network/bfban/checkban?names=" + id)
	res, err := cli.Request().Send()
	if err != nil {
		return false
	}
	return gjson.Get(res.String(), "names."+strings.ToLower(id)+".hacker").Bool()
}

//获取玩家pid
func GetPersonalID(name string) (string, error) {
	cli := gentleman.New()
	cli.URL("https://gateway.ea.com/proxy/identity/personas?namespaceName=cem_ea_id&displayName=" + name)
	cli.Use(headers.Set("X-Expand-Results", "true"))
	cli.Use(headers.Set("Authorization", rsp.TOKEN))
	cli.Use(headers.Set("Host", "gateway.ea.com"))
	res, err := cli.Request().Send()
	if err != nil {
		return "", err
	}
	info := gjson.Get(res.String(), "error").Str
	if info == "invalid_access_token" || info == "invalid_oauth_info" {
		rsp.Session(rsp.USERNAME, rsp.PASSWORD, true)
		return GetPersonalID(name)
	}
	if info != "" {
		return "", errors.New(info)
	}
	return gjson.Get(res.String(), "personas.persona.0.personaId").String(), err
}

//简体转繁体
func S2tw(str string) string {
	result := ""
	for _, v := range str {
		result += twmap[string(v)]
	}
	return result
}
