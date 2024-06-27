package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	bf1player "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/player"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/errcode"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/renderer"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/tracer"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/wdvxdr1123/ZeroBot/message"
	"go.opentelemetry.io/otel/codes"
)

// BindAccount 绑定账号
func (s *Service) BindAccount(ctx context.Context) error {
	_, span := global.Tracer.Start(ctx, "BindAccount")
	defer span.End()
	id := s.zctx.State["args"].(string)
	if id == "" {
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("绑定失败, ERR: 空id"))
		return errcode.InvalidPlayerError
	}
	// 数据库查询是否绑定
	span.AddEvent("start query database", tracer.AddEventWithDescription(tracer.Description("query player", id)))
	player, err := s.dao.GetPlayerByQID(s.zctx.Event.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("正在绑定id为 ", id))
		span.AddEvent("try to bind account")
		err = s.dao.CreatePlayer(s.zctx.Event.UserID, id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "database error")
			s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("绑定失败, ERR: 数据库错误"))
			return errcode.DataBaseCreateError
		}
		span.AddEvent("bind success")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("绑定成功"))
		return errcode.Success
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "database error")
		return errcode.DataBaseReadError.WithDetails("Error", err).WithZeroContext(s.zctx)
	}
	// 绑定的是旧id
	if id == player.DisplayName {
		span.AddEvent("detected user try to bind same, operation cancelled ")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("笨蛋! 你现在绑的就是这个id"))
		return errcode.Canceled
	}
	span.AddEvent("detected user try to change display name")
	s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("将原绑定id为 ", player.DisplayName, " 改绑为 ", id))
	err = s.dao.UpdatePlayer(s.zctx.Event.UserID, "", id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "database error")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("绑定失败, ERR: 数据库错误"))
		return errcode.DataBaseUpdateError
	}
	span.AddEvent("bind success")
	s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("绑定 ", id, " 成功"))
	span.AddEvent("try to get pid")
	pid, err := bf1api.GetPersonalID(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "network error")
		return errcode.NetworkError.WithDetails("bf1api.GetPersonalID", err).WithZeroContext(s.zctx)
	}
	span.AddEvent("write pid")
	err = s.dao.UpdatePlayer(s.zctx.Event.UserID, pid, "")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "database error")
		return errcode.DataBaseUpdateError
	}
	span.AddEvent("write success")
	span.SetStatus(codes.Ok, "")
	return errcode.Success
}

func (s *Service) getPlayerID(ctx context.Context) (string, bool) {
	_, span := global.Tracer.Start(ctx, "GetPlayerDisplayName")
	defer span.End()
	id := s.zctx.State["regex_matched"].([]string)[1]
	span.AddEvent("get name", tracer.AddEventWithDescription(tracer.Description("player name", id)))
	// id 为空就去数据库查
	if id == "" {
		span.AddEvent("empty name, query database")
		p, err := s.dao.GetPlayerByQID(s.zctx.Event.UserID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "database error")
			// 查不到或者失败就是无效
			return "", false
		}
		span.AddEvent("query success", tracer.AddEventWithDescription(tracer.Description("player name", p.DisplayName)))
		span.SetStatus(codes.Ok, "")
		return p.DisplayName, true
	}
	// id 不为空就认为有效
	span.SetStatus(codes.Ok, "")
	return id, true
}

func (s *Service) getPlayer(ctx context.Context, name string) (*model.Player, error) {
	_, span := global.Tracer.Start(ctx, "GetPlayer")
	defer span.End()
	if name == "" {
		span.AddEvent("detected name empty, query database")
		return s.dao.GetPlayerByQID(s.zctx.Event.UserID)
	}
	span.AddEvent("query database")
	player, err := s.dao.GetPlayerByName(name)
	if err == nil && player.PersonalID != "" {
		span.AddEvent("success")
		span.SetStatus(codes.Ok, "")
		return player, nil
	}
	span.AddEvent("detected pid empty")
	// 有错误或者pid残缺
	pid, err := bf1api.GetPersonalID(name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get pid failed")
		return nil, err
	}
	span.AddEvent("get pid success")
	span.SetStatus(codes.Ok, "")
	// 残缺
	if player != nil {
		player.PersonalID = pid
		span.AddEvent("rewrite pid")
		_ = player.Update(global.DB)
		return player, nil
	}
	return &model.Player{DisplayName: name, PersonalID: pid}, nil
}

func (s *Service) sendWeaponInfo(ctx context.Context, id, class string) error {
	nCtx, span := global.Tracer.Start(ctx, "SendWeaponInfo")
	defer span.End()
	span.AddEvent("try to get player")
	s.zctx.Send("少女折寿中...")
	player, err := s.getPlayer(nCtx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get player failed")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return errcode.NotFoundError.WithDetails("s.getPlayer", err).WithZeroContext(s.zctx)
	}
	span.AddEvent("try to get weapon info")
	weapons, err := bf1player.GetWeapons(player.PersonalID, class)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get weapon info failed")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 获取武器失败"))
		return errcode.NetworkError.WithDetails("bf1player.GetWeapons", err).WithZeroContext(s.zctx)
	}

	txt := "id：" + player.DisplayName + "\n"
	wp := ([]bf1player.Weapons)(*weapons)
	for i := 0; i < 5; i++ {
		txt += fmt.Sprintf("%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n",
			"---------------",
			"武器名：", wp[i].Name,
			"击杀数：", strconv.FormatFloat(wp[i].Kills, 'f', 0, 64),
			"准度：", wp[i].Accuracy,
			"爆头率：", wp[i].Headshots,
			"KPM：", wp[i].KPM,
			"效率：", wp[i].Efficiency,
		)
	}

	span.AddEvent("success")
	span.SetStatus(codes.Ok, "")
	renderer.Txt2Img(s.zctx, txt)
	return errcode.Success
}

// GetPlayerRecent 获取玩家最近游玩
func (s *Service) GetPlayerRecent(ctx context.Context) error {
	nCtx, span := global.Tracer.Start(ctx, "GetPlayerRecent")
	defer span.End()
	span.AddEvent("try to get player name")
	id, isVaild := s.getPlayerID(nCtx)
	if !isVaild {
		span.AddEvent("failed to get player")
		span.SetStatus(codes.Error, "get invalid player display name")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return errcode.NotFoundError
	}
	s.zctx.Send("少女折寿中...")
	span.AddEvent("try to get player recent")
	recent, err := bf1player.GetBF1Recent(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "network error")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 获取最近战绩失败"))
		return errcode.NetworkError.WithDetails("bf1player.GetBF1Recent", err).WithZeroContext(s.zctx)
	}
	// 发送最近战绩
	// TODO: 修改为卡片发送
	msg := "id：" + id + "\n"
	for i := range *recent {
		msg += "服务器：" + (*recent)[i].Server[:24] + "\n"
		msg += "地图：" + (*recent)[i].Map + "   (" + (*recent)[i].Mode + ")\n"
		msg += "kd：" + strconv.FormatFloat((*recent)[i].Kd, 'f', -1, 64) + "\n"
		msg += "kpm：" + strconv.FormatFloat((*recent)[i].Kpm, 'f', -1, 64) + "\n"
		msg += "游玩时长：" + strconv.FormatFloat(float64((*recent)[i].Time/60), 'f', -1, 64) + "分钟"
		msg += "\n---------------\n"
	}
	span.AddEvent("success")
	span.SetStatus(codes.Ok, "")
	renderer.Txt2Img(s.zctx, msg)
	return errcode.Success
}

// GetPlayerStats 获取玩家战绩
func (s *Service) GetPlayerStats(ctx context.Context) error {
	nCtx, span := global.Tracer.Start(ctx, "GetPlayerStats")
	defer span.End()
	span.AddEvent("try to get player name")
	id, isVaild := s.getPlayerID(nCtx)
	if !isVaild {
		span.AddEvent("failed to get player")
		span.SetStatus(codes.Error, "get invalid player display name")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return errcode.NotFoundError
	}
	s.zctx.Send("少女折寿中...")
	span.AddEvent("try to get player status")
	stat, err := bf1player.GetStats(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "network error")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 获取玩家战绩失败, 请自行检查id是否正确"))
		return errcode.NetworkError.WithDetails("bf1player.GetStats", err).WithZeroContext(s.zctx)
	}
	if stat.Rank == "" {
		span.SetStatus(codes.Error, "invalid status")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("获取到的部分数据为空，请检查id是否有效"))
		return errcode.InternalError.WithDetails("Unexpected", errors.Errorf("%s stat.Rank is blank", id))
	}
	// 发送战绩
	// TODO: 修改为卡片发送, 部分数据不准确，等待更改
	txt := "id：" + id + "\n" +
		"等级：" + stat.Rank + "\n" +
		"游玩时长：" + stat.TimePlayed + "\n" +
		"总kd：" + stat.TotalKD + "(" + stat.Kills + "/" + stat.Deaths + ")" + "\n" +
		"总kpm：" + stat.KPM + "\n" +
		"准度：" + stat.Accuracy + "\n" +
		"爆头数：" + stat.Headshots + "\n" +
		"胜率：" + stat.WinPercent + "(" + stat.Wins + "/" + stat.Losses + ")" + "\n" +
		"场均击杀：" + stat.KillsPerGame + "\n" +
		"步战kd：" + stat.InfantryKD + "\n" +
		"步战击杀：" + stat.InfantryKills + "\n" +
		"步战kpm：" + stat.InfantryKPM + "\n" +
		"载具击杀：" + stat.VehicleKills + "\n" +
		"载具kpm：" + stat.VehicleKPM + "\n" +
		"近战击杀：" + stat.DogtagsTaken + "\n" +
		"最高连杀：" + stat.HighestKillStreak + "\n" +
		"最远爆头：" + stat.LongestHeadshot + "\n" +
		"MVP数：" + stat.MVP + "\n" +
		"作为神医拉起了 " + stat.Revives + " 人"

	span.AddEvent("success")
	span.SetStatus(codes.Ok, "")
	renderer.Txt2Img(s.zctx, txt)
	return errcode.Success
}

// GetPlayerWeapon 获取玩家武器
func (s *Service) GetPlayerWeapon(ctx context.Context) error {
	nCtx, span := global.Tracer.Start(ctx, "GetPlayerWeapon")
	defer span.End()

	span.AddEvent("phase user command")
	str := strings.Split(s.zctx.State["regex_matched"].([]string)[1], " ")
	var id string
	span.SetStatus(codes.Ok, "")
	// 相当于只输入 .武器
	if str[0] == "" {
		return s.sendWeaponInfo(nCtx, id, bf1player.ALL)
	}
	if len(str) > 1 {
		id = str[1]
	}
	switch str[0] {
	// 除default 相当于输入 .武器 class id
	case "半自动", "semi":
		return s.sendWeaponInfo(nCtx, id, bf1player.Semi)
	case "冲锋枪", "冲锋":
		return s.sendWeaponInfo(nCtx, id, bf1player.SMG)
	case "轻机枪", "机枪":
		return s.sendWeaponInfo(nCtx, id, bf1player.LMG)
	case "步枪", "狙击枪", "狙击":
		return s.sendWeaponInfo(nCtx, id, bf1player.Bolt)
	case "霰弹枪", "散弹枪", "霰弹", "散弹":
		return s.sendWeaponInfo(nCtx, id, bf1player.Shotgun)
	case "配枪", "手枪", "副手":
		return s.sendWeaponInfo(nCtx, id, bf1player.Sidearm)
	case "近战", "刀":
		return s.sendWeaponInfo(nCtx, id, bf1player.Melee)
	case "手榴弹", "手雷", "雷":
		return s.sendWeaponInfo(nCtx, id, bf1player.Grenade)
	case "驾驶员", "坦克兵", "载具":
		return s.sendWeaponInfo(nCtx, id, bf1player.Dirver)
	case "配备", "装备":
		return s.sendWeaponInfo(nCtx, id, bf1player.Gadget)
	case "精英", "精英兵":
		return s.sendWeaponInfo(nCtx, id, bf1player.Elite)
	default:
		// 相当于 .武器 id
		if regexp.MustCompile(`\w+`).MatchString(str[0]) {
			id = str[0]
			return s.sendWeaponInfo(nCtx, id, bf1player.ALL)
		}
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 获取玩家武器失败, 不能识别的输入格式"))
		return errcode.InvalidParamsError.WithZeroContext(s.zctx)
	}
}

// GetPlayerVehicle 获取玩家载具信息
func (s *Service) GetPlayerVehicle(ctx context.Context) error {
	nCtx, span := global.Tracer.Start(ctx, "GetPlayerVehicle")
	defer span.End()
	s.zctx.Send("少女折寿中...")
	span.AddEvent("try to get player")
	player, err := s.getPlayer(nCtx, s.zctx.State["regex_matched"].([]string)[1])
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get player failed")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return errcode.NotFoundError
	}
	span.AddEvent("try to get vehicle")
	car, err := bf1player.GetVehicles(player.PersonalID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get vehicle failed")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 获取玩家载具失败, 请自行检查id是否正确"))
		return errcode.NetworkError.WithDetails("bf1player.GetVehicles", err).WithZeroContext(s.zctx)
	}
	msg := "id：" + player.DisplayName + "\n"
	for i := range *car {
		msg += "------------\n"
		msg += (*car)[i].Name + "\n"
		msg += fmt.Sprintf("%s%6.0f\t", "击杀数：", (*car)[i].Kills)
		msg += "kpm：" + (*car)[i].KPM + "\n"
		msg += fmt.Sprintf("%s%6.0f\t", "击毁数：", (*car)[i].Destroyed)
		msg += "游玩时间：" + (*car)[i].Time + " 小时\n"
	}
	span.AddEvent("success")
	span.SetStatus(codes.Ok, "")
	renderer.Txt2Img(s.zctx, msg)
	return errcode.Success
}

// GetBF1Exchange 获取BF1本期交换信息
func (s *Service) GetBF1Exchange(ctx context.Context) error {
	_, span := global.Tracer.Start(ctx, "GetBF1Exchange")
	defer span.End()
	span.AddEvent("try to get exchange")
	exchange, err := bf1api.GetExchange()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get exchange failed")
		s.zctx.SendChain(message.Reply(s.zctx.Event.MessageID), message.Text("ERR: 获取交换失败"))
		return errcode.NetworkError.WithDetails("bf1api.GetExchange", err).WithZeroContext(s.zctx)
	}
	var msg string
	for i, v := range exchange {
		msg += i + ": \n"
		for _, skin := range v {
			msg += "\t" + skin + "\n"
		}
	}
	span.AddEvent("success")
	span.SetStatus(codes.Ok, "")
	renderer.Txt2Img(s.zctx, msg)
	return errcode.Success
}

// GetBF1OpreationPack 获取本期行动包信息
func (s *Service) GetBF1OpreationPack(ctx context.Context) error {
	_, span := global.Tracer.Start(ctx, "GetBF1OpreationPack")
	defer span.End()
	span.AddEvent("try to get pack")
	pack, err := bf1api.GetCampaignPacks()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get pack failed")
		s.zctx.SendChain(message.Reply(s.zctx.Event.MessageID), message.Text("ERR: 获取行动包失败"))
		return errcode.NetworkError.WithDetails("bf1api.GetCampaignPacks", err).WithZeroContext(s.zctx)
	}
	var msg string
	msg += "行动名：" + pack.Name + "\n"
	msg += "剩余时间：" + fmt.Sprintf("%.2f", float64(pack.RemainTime)/60/24) + " 天\n"
	msg += "箱子重置时间：" + fmt.Sprintf("%.2f", float64(pack.ResetTime)/60) + " 小时\n"
	msg += "行动地图：" + pack.Op1Name + " 与 " + pack.Op2Name + "\n"
	msg += "行动简介：" + pack.Desc

	span.AddEvent("success")
	span.SetStatus(codes.Ok, "")
	renderer.Txt2Img(s.zctx, msg)
	return errcode.Success
}

// GetPlayerBanInfo 获取玩家联ban信息
func (s *Service) GetPlayerBanInfo(ctx context.Context) error {
	nCtx, span := global.Tracer.Start(ctx, "GetPlayerBanInfo")
	defer span.End()
	s.zctx.Send("少女折寿中...")
	span.AddEvent("try to get player name")
	id, isVaild := s.getPlayerID(nCtx)
	if !isVaild {
		span.AddEvent("failed to get player")
		span.SetStatus(codes.Error, "get invalid player display name")
		s.zctx.SendChain(message.At(s.zctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return errcode.NotFoundError
	}
	span.AddEvent("try to get player ban info")
	info := bf1player.IsHacker(id)
	var msg string
	msg += "id: " + id + "\n"
	msg += "EAC: " + "\n\t" + info.EAC.Status + "\n\t案件链接: " + info.EAC.URL + "\n"
	msg += "BFBan: " + "\n\t状态: " + info.BFBan.Status + "\n\t"
	if info.BFBan.IsCheater {
		msg += "案件链接: " + info.BFBan.URL
	}
	span.AddEvent("success")
	span.SetStatus(codes.Ok, "")
	s.zctx.SendChain(message.Reply(s.zctx.Event.MessageID), message.Text(msg))
	return errcode.Success
}
