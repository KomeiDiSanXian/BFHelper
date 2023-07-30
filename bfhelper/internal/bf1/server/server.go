// Package server 战地1服务器操作
package server

import (
	"errors"
	"fmt"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	bf1reqbody "github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq/bf1"
)

// Server 服务器结构体
type Server struct {
	SID  string
	GID  string
	PGID string
}

type m struct {
	MapName  string
	ModeName string
}

// Maps 地图切片
type Maps []m

// NewServer 服务器
func NewServer(sid, gid, pgid string) *Server {
	return &Server{
		SID:  sid,
		GID:  gid,
		PGID: pgid,
	}
}

// Kick player, reason needs BIG5, return reason and err
func (s *Server) Kick(pid, reason string) (string, error) {
	reason = fmt.Sprintf("%s%s", "Remi:", reason)
	if len(reason) > 32 {
		return "", errors.New("理由过长")
	}
	post := bf1reqbody.NewPostKick(pid, s.GID, reason)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return "", err
	}
	return data.Get("result.reason").Str, err
}

// Ban player, check returned id
func (s *Server) Ban(pid string) error {
	post := bf1reqbody.NewPostBan(pid, s.SID)
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
func (s *Server) Unban(pid string) error {
	post := bf1reqbody.NewPostRemoveBan(pid, s.SID)
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
func (s *Server) ChangeMap(index int) error {
	post := bf1reqbody.NewPostChangeMap(s.PGID, index)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if data.Get("id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// GetMaps returns maps
func (s *Server) GetMaps() (*Maps, error) {
	post := bf1reqbody.NewPostGetServerInfo(s.GID)
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
	mp := make(Maps, len(result))
	for i, v := range result {
		mp[i] = m{MapName: v.Get("mapPrettyName").Str, ModeName: v.Get("modePrettyName").Str}
	}
	return &mp, nil
}

// GetAdminspid returns pids of admins
func (s *Server) GetAdminspid() ([]string, error) {
	post := bf1reqbody.NewPostRSPInfo(s.SID)
	data, err := bf1api.ReturnJSON(global.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	result := data.Get("result.adminList.#.personaId").Array()
	result = append(result, data.Get("result.owner.personaId"))
	strs := make([]string, len(result))
	for i, v := range result {
		strs[i] = v.Str
	}
	return strs, bf1api.Exception(data.Get("error.code").Int())
}

// input keywords for map id
/* not compiled
func (s *Server) GetMapidByKeywords(keyword string) (int, error) {
	switch keyword{
		case
	}
}
*/
