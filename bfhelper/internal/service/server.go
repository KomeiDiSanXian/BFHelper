package service

import (
	"strconv"
	"strings"
	"sync"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	bf1server "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/server"
	"github.com/pkg/errors"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// CreateGroup 所在群创建一个服务器群组
//
// @permission: GroupOwner
func (s *Service) CreateGroup() error {
	owner := s.ctx.State["args"].(int64)
	create := func(o int64) error {
		err := s.dao.CreateGroup(s.ctx.Event.GroupID, o)
		if err != nil {
			s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 创建服务器群组失败"))
			return err
		}
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("创建成功!"))
		return nil
	}
	s.ctx.Send("少女折寿中...")
	if owner != 0 {
		return create(owner)
	}
	return create(s.ctx.Event.UserID)
}

// DeleteGroup 所在群删除服务器群组
//
// @permission: ServerOwner
func (s *Service) DeleteGroup() error {
	err := s.dao.DeleteGroup(s.ctx.Event.GroupID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 删除服务器群组失败"))
		return err
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("删除成功"))
	return nil
}

// ChangeOwner 更换服务器群组所有人
//
// @permission: ServerOwner
func (s *Service) ChangeOwner() error {
	owner := s.ctx.State["args"].(int64)
	if owner == 0 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 更换服务器群组所属失败"))
		return errors.New("invalid owner")
	}
	err := s.dao.UpdateOwner(s.ctx.Event.GroupID, owner)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 更换服务器群组所属失败"))
		return err
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("更换成功"))
	return nil
}

func (s *Service) getServer(gameID string) (*bf1server.Server, error) {
	result, err := bf1api.GetServerFullInfo(gameID)
	if err != nil {
		return nil, err
	}
	server := bf1server.NewServer(
		result.Get("result.rspInfo.server.serverId").Str,
		result.Get("result.serverInfo.gameId").Str,
		result.Get("result.serverInfo.guid").Str,
	)
	server.Name = result.Get("result.serverInfo.name").Str
	return server, nil
}

func (s *Service) addServerProcess(gameID string, groupID int64) error {
	server, err := s.getServer(gameID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 添加服务器 ", gameID, " 失败"))
		return err
	}
	err = s.dao.AddGroupServer(groupID, server.GID, server.SID, server.PGID, server.Name)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 添加服务器 ", gameID, " 失败"))
		return err
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("成功添加服务器 ", server.GID))
	return nil
}

// AddServer 添加服务器
//
// @permission: SuperAdmin
func (s *Service) AddServer() error {
	str := s.ctx.State["args"].(string)
	if str == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 无效的输入"))
		return errors.New("invalid input")
	}
	strs := strings.Split(str, " ")
	groupID, _ := strconv.ParseInt(strs[0], 10, 64)
	if len(strs) < 2 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: GameID 为空"))
		return errors.New("invalid gameid")
	}
	var wg sync.WaitGroup
	for i := 1; i < len(strs); i++ {
		go func(gameID string) {
			wg.Add(1)
			_ = s.addServerProcess(gameID, groupID)
			wg.Done()
		}(strs[i])
	}
	wg.Wait()
	return nil
}

// AddServerAdmin 添加服务器管理员
//
// @permission: ServerOwner
func (s *Service) AddServerAdmin() error {
	str := s.ctx.State["args"].(string)
	if str == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 填写的管理员qq为空"))
		return errors.New("invalid admin")
	}
	admins := strings.Split(str, " ")
	qqs := make([]int64, 0, len(admins))
	for _, a := range admins {
		adminQQ, _ := strconv.ParseInt(a, 10, 64)
		qqs = append(qqs, adminQQ)
	}
	err := s.dao.AddGroupAdmin(s.ctx.Event.GroupID, qqs...)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 添加管理失败"))
		return err
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("添加成功"))
	return nil
}

// SetServerAlias 设置服务器别名
//
// @permission: ServerAdmin
func (s *Service) SetServerAlias() error {
	str := s.ctx.State["args"].(string)
	if str == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errors.New("invalid input")
	}
	strs := strings.Split(str, " ")
	if len(strs) != 2 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 不能识别的输入"))
		return errors.New("invalid input")
	}
	gameID, alias := strs[0], strs[1]
	err := s.dao.ServerAlias(gameID, alias)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 设置别名失败"))
		return err
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("设置成功"))
	return nil
}

// DeleteServer 删除服务器
//
// @permission: ServerOwner
func (s *Service) DeleteServer() error {
	gameID := s.ctx.State["args"].(string)
	if gameID == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errors.New("invalid input")
	}
	err := s.dao.RemoveGroupServer(s.ctx.Event.GroupID, gameID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 删除服务器 ", gameID, " 失败"))
		return err
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("删除成功"))
	return nil
}

// DeleteAdmin 删除群组服务器管理员
//
// @permission: ServerOwner
func (s *Service) DeleteAdmin() error {
	qq := s.ctx.State["args"].(int64)
	if qq == 0 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errors.New("invalid input")
	}
	err := s.dao.RemoveGroupAdmin(s.ctx.Event.GroupID, qq)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 删除管理员 ", qq, " 失败"))
		return err
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("删除成功"))
	return nil
}
