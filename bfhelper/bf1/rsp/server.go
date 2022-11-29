package bf1rsp

import (
	"errors"
	"fmt"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	"github.com/tidwall/gjson"
)

type server struct {
	Sid  string
	Gid  string
	PGid string
}

type m struct {
	MapName  string
	ModeName string
}

type maps []m

// New server info
func NewServer(sid, gid, pgid string) *server {
	return &server{
		Sid:  sid,
		Gid:  gid,
		PGid: pgid,
	}
}

// kick player, reason needs BIG5, return reason and err
func (s *server) Kick(pid, reason string) (string, error) {
	reason = fmt.Sprintf("%s%s", "Remi:", reason)
	if len(reason) > 32 {
		return "", errors.New("理由过长")
	}
	post := NewPostKick(pid, s.Gid, reason)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return "", err
	}
	return gjson.Get(data, "result.reason").Str, err
}

// ban player, check returned id
func (s *server) Ban(pid string) error {
	post := NewPostBan(pid, s.Sid)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// unban player
func (s *server) Unban(pid string) error {
	post := NewPostRemoveBan(pid, s.Sid)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// choose level
func (s *server) ChangeMap(index int) error {
	post := NewPostChangeMap(s.PGid, index)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// get level
func (s *server) GetMaps() (*maps, error) {
	post := NewPostGetServerInfo(s.Gid)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	if gjson.Get(data, "result").String() == "" {
		return nil, errors.New("服务器gameid可能无效，请更新服务器信息")
	}
	result := gjson.Get(data, "result.rotation").Array()
	if result == nil {
		return nil, errors.New("获取到的地图池为空")
	}
	var mp maps
	for _, v := range result {
		mp = append(mp, m{MapName: v.Get("mapPrettyName").Str, ModeName: v.Get("modePrettyName").Str})
	}
	return &mp, nil
}

func (s *server) GetAdminspid() ([]string, error) {
	post := NewPostRSPInfo(s.Sid)
	data, err := bf1api.ReturnJson(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	result := gjson.Get(data, "result.adminList.#.personaId").Array()
	result = append(result, gjson.Get(data, "result.owner.personaId"))
	var strs []string
	for _, v := range result {
		strs = append(strs, v.Str)
	}
	return strs, bf1api.Exception(gjson.Get(data, "error.code").Int())
}
