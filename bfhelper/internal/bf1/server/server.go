// Package server 战地1服务器操作
package server

import (
	"errors"
	"fmt"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	bf1reqbody "github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq/bf1"
	"github.com/tidwall/gjson"
)

// Map 服务器地图
type Map struct {
	Name   string
	Mode   string
}

// Kick player, reason needs BIG5, return reason and err
func Kick(gameID, pid, reason string) (string, error) {
	reason = fmt.Sprintf("%s%s", "Remi:", reason)
	if len(reason) > 32 {
		return "", errors.New("理由过长")
	}
	post := bf1reqbody.NewPostKick(pid, gameID, reason)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return "", err
	}
	return data.Get("result.reason").Str, err
}

// Ban player, check returned id
func Ban(serverID, pid string) error {
	post := bf1reqbody.NewPostBan(pid, serverID)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if data.Get("id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// Unban player
func Unban(serverID, pid string) error {
	post := bf1reqbody.NewPostRemoveBan(pid, serverID)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if data.Get("id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// ChangeMap will change the map for players
func ChangeMap(pgid string, index int) error {
	post := bf1reqbody.NewPostChangeMap(pgid, index)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if data.Get("id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

func mapRequest(gameID string) ([]gjson.Result, error) {
	post := bf1reqbody.NewPostGetServerInfo(gameID)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	if data.Get("result").String() == "" {
		return nil, errors.New("服务器gameid可能无效，请更新服务器信息")
	}
	result := data.Get("result.rotation").Array()
	if result == nil {
		return nil, errors.New("获取到的地图池为空")
	}
	return result, nil
}

// GetMapSlice returns map slice
func GetMapSlice(gameID string) ([]*Map, error) {
	result, err := mapRequest(gameID)
	if err != nil {
		return nil, err
	}

	mp := make([]*Map, 0, len(result))
	for _, v := range result {
		m := &Map{Name: v.Get("mapPrettyName").Str, Mode: v.Get("modePrettyName").Str}
		mp = append(mp, m)
	}
	return mp, nil
}
