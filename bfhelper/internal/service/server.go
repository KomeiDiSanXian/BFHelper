package service

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	bf1server "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/server"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/errcode"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/textutil"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/renderer"	
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type server struct {
	name, gameID, serverID, pgid string
}

// CreateGroup 所在群创建一个服务器群组
//
// @permission: GroupOwner
func (s *Service) CreateGroup() error {
	o := s.ctx.State["args"].(string)
	owner, _ := strconv.ParseInt(o, 10, 64)
	create := func(o int64) error {
		err := s.dao.CreateGroup(s.ctx.Event.GroupID, o)
		if err != nil {
			s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 创建服务器群组失败"))
			return errcode.DataBaseCreateError
		}
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("创建成功!"))
		return errcode.Success
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
		return errcode.DataBaseDeleteError
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("删除成功"))
	return errcode.Success
}

// ChangeOwner 更换服务器群组所有人
//
// @permission: ServerOwner
func (s *Service) ChangeOwner() error {
	o := s.ctx.State["args"].(string)
	owner, _ := strconv.ParseInt(o, 10, 64)
	if owner == 0 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 未输入要更换的QQ号"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	err := s.dao.UpdateOwner(s.ctx.Event.GroupID, owner)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 更换服务器群组所属失败"))
		return errcode.DataBaseUpdateError
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("更换成功"))
	return errcode.Success
}

func (s *Service) getServer(gameID string) (*server, error) {
	result, err := bf1api.GetServerFullInfo(gameID)
	if err != nil {
		return nil, err
	}
	srv := server{
		serverID: result.Get("result.rspInfo.server.serverId").Str,
		gameID:   result.Get("result.serverInfo.gameId").Str,
		pgid:     result.Get("result.serverInfo.guid").Str,
	}
	srv.name = result.Get("result.serverInfo.name").Str
	return &srv, nil
}

func (s *Service) addServerProcess(gameID string, groupID int64) error {
	server, err := s.getServer(gameID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 添加服务器 ", gameID, " 失败"))
		return err
	}
	err = s.dao.AddGroupServer(groupID, server.gameID, server.serverID, server.pgid, server.name)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 添加服务器 ", gameID, " 失败"))
		return err
	}
	return nil
}

func (s *Service) addServer(gameID string, groupID int64, wg *sync.WaitGroup, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	defer wg.Done()
	_ = s.addServerProcess(gameID, groupID)
}

func (s *Service) addServers(gameIDs []string, groupID int64) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, gameID := range gameIDs {
		wg.Add(1)
		go s.addServer(gameID, groupID, &wg, &mu)
	}
	wg.Wait()
}

// AddServer 添加服务器
//
// @permission: SuperAdmin
func (s *Service) AddServer() error {
	str := s.ctx.State["args"].(string)
	if str == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 无效的输入"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	strs := strings.Split(str, " ")
	groupID, _ := strconv.ParseInt(strs[0], 10, 64)
	if len(strs) < 2 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: GameID 为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	s.addServers(strs[1:], groupID)
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("添加完成"))
	return errcode.Success
}

// AddServerAdmin 添加服务器管理员
//
// @permission: ServerOwner
func (s *Service) AddServerAdmin() error {
	str := s.ctx.State["args"].(string)
	if str == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 填写的管理员qq为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
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
		return errcode.DataBaseUpdateError
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("添加成功"))
	return errcode.Success
}

// SetServerAlias 设置服务器别名
//
// @permission: ServerAdmin
func (s *Service) SetServerAlias() error {
	str := s.ctx.State["args"].(string)
	if str == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	strs := strings.Split(str, " ")
	if len(strs) != 2 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 不能识别的输入"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	gameID, alias := strs[0], strs[1]
	err := s.dao.ServerAlias(gameID, alias)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 设置别名失败"))
		return errcode.DataBaseUpdateError
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("设置成功"))
	return errcode.Success
}

// DeleteServer 删除服务器
//
// @permission: ServerOwner
func (s *Service) DeleteServer() error {
	gameID := s.ctx.State["args"].(string)
	if gameID == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	err := s.dao.RemoveGroupServer(s.ctx.Event.GroupID, gameID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 删除服务器 ", gameID, " 失败"))
		return errcode.DataBaseDeleteError
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("删除成功"))
	return errcode.Success
}

// DeleteAdmin 删除群组服务器管理员
//
// @permission: ServerOwner
func (s *Service) DeleteAdmin() error {
	o := s.ctx.State["args"].(string)
	qq, _ := strconv.ParseInt(o, 10, 64)
	if qq == 0 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	err := s.dao.RemoveGroupAdmin(s.ctx.Event.GroupID, qq)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 删除管理员 ", qq, " 失败"))
		return errcode.DataBaseDeleteError
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("删除成功"))
	return errcode.Success
}

func (s *Service) kickProcess(server model.Server, pid, reason string, msgChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	name := server.NameInGroup
	if name == "" {
		name = server.ServerName
	}
	returned, err := bf1server.Kick(server.GameID, pid, reason)
	if err != nil {
		msgChan <- fmt.Sprintf("ERROR: 在 %s%s\n", name, " 踢出失败")
	} else {
		msgChan <- fmt.Sprintf("在 %s%s%s\n", name, " 踢出成功: ", returned)
	}
}

func (s *Service) kick(pid string, reason string, group *model.Group) {
	var wg sync.WaitGroup
	var tosend string
	msgChan := make(chan string, len(group.Servers))
	for _, server := range group.Servers {
		wg.Add(1)
		go s.kickProcess(server, pid, reason, msgChan, &wg)
	}
	wg.Wait()
	close(msgChan)
	for msg := range msgChan {
		tosend += msg
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text(tosend))
}

// KickPlayer 在绑定的服务器里踢出玩家
//
// @permission: ServerAdmin
func (s *Service) KickPlayer() error {
	cmdString := s.ctx.State["args"].(string)
	if cmdString == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	cmds := strings.Split(cmdString, " ")
	group, err := s.dao.GetGroup(s.ctx.Event.GroupID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 获取绑定服务器失败"))
		return errcode.DataBaseReadError.WithDetails("Error", err).WithZeroContext(s.ctx)
	}
	player := cmds[0]
	reason := "Admin kick"
	if len(cmds) >= 2 {
		reason = cmds[1]
	}
	reason = textutil.Traditionalize(reason)
	if cleaned, has := textutil.CleanPersonalID(player); has {
		s.kick(cleaned, reason, group)
		return errcode.Success
	}
	pl, err := s.getPlayer(player)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 踢出失败: 未查询到目标玩家pid"))
		return errcode.NotFoundError.WithDetails("Error", err).WithZeroContext(s.ctx)
	}
	s.kick(pl.PersonalID, reason, group)
	return errcode.Success
}

// 单服务器封禁/解封
func (s *Service) banFunc(banfunc func(sid string, pid string) error) error {
	cmdString := s.ctx.State["args"].(string)
	if cmdString == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	cmds := strings.Split(cmdString, " ")
	if len(cmds) != 2 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 不能识别的输入"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	srvName, player := cmds[0], cmds[1]
	srv, err := s.dao.GetServerByAlias(s.ctx.Event.GroupID, srvName)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 查询服务器失败"))
		return errcode.DataBaseReadError.WithDetails("Error", err).WithZeroContext(s.ctx)
	}
	if cleaned, has := textutil.CleanPersonalID(player); has {
		err := banfunc(srv.ServerID, cleaned)
		if err != nil {
			s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 失败"))
			return errcode.InternalError.WithDetails("Error", err).WithZeroContext(s.ctx)
		}
		return errcode.Success
	}
	pl, err := s.getPlayer(player)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 未查询到目标玩家pid"))
		return errcode.NotFoundError.WithDetails("Error", err)
	}
	err = banfunc(srv.ServerID, pl.PersonalID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 失败"))
		return errcode.InternalError.WithDetails("Error", err).WithZeroContext(s.ctx)
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("操作", pl.DisplayName, "成功"))
	return errcode.Success
}

// 多服务器封禁/解封
func (s *Service) bansFunc(banfunc func(sid string, pid string) error) error {
	playerName := s.ctx.State["args"].(string)
	if playerName == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	player, err := s.getPlayer(playerName)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 未查询到目标玩家pid"))
		return errcode.NotFoundError.WithDetails("Error", err)
	}
	group, err := s.dao.GetGroup(s.ctx.Event.GroupID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 获取绑定服务器失败"))
		return errcode.DataBaseReadError.WithDetails("Error", err).WithZeroContext(s.ctx)
	}

	var wg sync.WaitGroup
	var tosend string
	msgChan := make(chan string, len(group.Servers))
	for _, server := range group.Servers {
		wg.Add(1)
		go s.bansProcess(&server, player, msgChan, &wg, banfunc)
	}
	wg.Wait()
	close(msgChan)
	for msg := range msgChan {
		tosend += msg
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text(tosend))
	return errcode.Success
}

func (s *Service) bansProcess(srv *model.Server, player *model.Player, msgChan chan string, wg *sync.WaitGroup, banfunc func(sid string, pid string) error) {
	defer wg.Done()
	srvName := srv.NameInGroup
	if srvName == "" {
		srvName = srv.ServerName
	}
	err := banfunc(srv.ServerID, player.PersonalID)
	if err != nil {
		msgChan <- fmt.Sprintf("ERROR: 在 %s%s\n", srvName, " 操作失败")
		return
	}
	msgChan <- fmt.Sprintf("在 %s%s\n", srvName, " 操作成功")
}

// BanPlayer 指定一个已绑定的服务器中封禁玩家
//
// @permission: ServerAdmin
// TODO: 添加 vban
func (s *Service) BanPlayer() error {
	return s.banFunc(bf1server.Ban)
}

// UnbanPlayer 指定一个已绑定的服务器中解封玩家
//
// @permission: ServerAdmin
func (s *Service) UnbanPlayer() error {
	return s.banFunc(bf1server.Unban)
}

// BanPlayerAtAllServer 在所有已绑定的服务器里封禁玩家
//
// @permission: ServerAdmin
func (s *Service) BanPlayerAtAllServer() error {
	return s.bansFunc(bf1server.Ban)
}

// UnbanPlayerAtAllServer 在所有已绑定的服务器里封禁玩家
//
// @permission: ServerAdmin
func (s *Service) UnbanPlayerAtAllServer() error {
	return s.bansFunc(bf1server.Unban)
}

func (s *Service) sendMaps(maptxt string, next *zero.FutureEvent, srv *model.Server, maps []*bf1server.Map) error {
	renderer.Txt2Img(s.ctx, maptxt)
	recv, cancle := next.Repeat()
	defer cancle()
	tick := time.NewTimer(time.Minute)
	for {
		select {
		case <-tick.C:
			s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 等待回复超时"))
			return errcode.TimeoutError.WithZeroContext(s.ctx)
		case c := <-recv:
			index, _ := strconv.Atoi(c.Event.Message.String())
			if index >= len(maps) || index < 0 {
				s.ctx.SendChain(message.Reply(c.Event.MessageID), message.Text("ERROR: 无效的地图序号,取值范围为 0-", len(maps)-1))
				return errcode.InvalidParamsError.WithDetails("MapID", "out of index").WithZeroContext(c)
			}
			err := bf1server.ChangeMap(srv.PGID, index)
			if err != nil {
				s.ctx.SendChain(message.Reply(c.Event.MessageID), message.Text("ERROR: 切图失败"))
				return errcode.NetworkError.WithDetails("bf1server.ChangeMap", err)
			}
			s.ctx.SendChain(message.Reply(c.Event.MessageID), message.Text("已切到 ", maps[index].Name, "(", maps[index].Mode, ")"))
			return errcode.Success
		}
	}
}

// ChangeMap 切换指定服务器的地图
//
// @permission: ServerAdmin
func (s *Service) ChangeMap() error {
	cmdString := s.ctx.State["args"].(string)
	if cmdString == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	cmds := strings.Split(cmdString, " ")
	srvName := cmds[0]
	srv, err := s.dao.GetServerByAlias(s.ctx.Event.GroupID, srvName)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 找不到别名为 ", srvName, " 的服务器"))
		return errcode.NotFoundError.WithDetails("Error", err).WithZeroContext(s.ctx)
	}
	maps, err := bf1server.GetMapSlice(srv.GameID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 获取地图池失败"))
		return errcode.NetworkError.WithDetails("Error", err)
	}

	maptxt := "请在一分钟内选择一个序号来回复\n------\n图池序号和模式\n"
	for i, m := range maps {
		maptxt += fmt.Sprintf("\t%2d %s(%s)\n", i, m.Name, m.Mode)
	}

	next := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^\d{1,2}$`), zero.OnlyGroup, s.ctx.CheckSession())
	if len(cmds) == 1 {
		return s.sendMaps(maptxt, next, srv, maps)
	}

	index, _ := strconv.Atoi(cmds[1])
	if index >= len(maps) || index < 0 {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 无效的地图序号,取值范围为 0-", len(maps)-1))
		return errcode.InvalidParamsError.WithDetails("MapID", "out of index").WithZeroContext(s.ctx)
	}
	err = bf1server.ChangeMap(srv.PGID, index)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 切图失败"))
		return errcode.NetworkError.WithDetails("bf1server.ChangeMap", err)
	}
	s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("已切到 ", maps[index].Name, "(", maps[index].Mode, ")"))
	return errcode.Success
}

// GetMap 查看地图池
//
// @permission: Everyone
func (s *Service) GetMap() error {
	srvName := s.ctx.State["args"].(string)
	if srvName == "" {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 输入为空"))
		return errcode.InvalidParamsError.WithZeroContext(s.ctx)
	}
	srv, err := s.dao.GetServerByAlias(s.ctx.Event.GroupID, srvName)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 找不到别名为 ", srvName, " 的服务器"))
		return errcode.NotFoundError.WithDetails("Error", err).WithZeroContext(s.ctx)
	}
	maps, err := bf1server.GetMapSlice(srv.GameID)
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERROR: 获取地图池失败"))
		return errcode.NetworkError.WithDetails("bf1server.GetMapSlice", err)
	}
	maptxt := "图池序号和模式\n"
	for i, m := range maps {
		maptxt += fmt.Sprintf("\t%2d %s(%s)\n", i, m.Name, m.Mode)
	}
	renderer.Txt2Img(s.ctx, maptxt)
	return errcode.Success
}
