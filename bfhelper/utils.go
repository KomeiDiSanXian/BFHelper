// Package bfhelper 工具
package bfhelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/jinzhu/gorm"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	api "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/player"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/netreq"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/renderer"
)

// ReturnBindID 检查是否绑定，返回id
func ReturnBindID(ctx *zero.Ctx, id string) (string, error) {
	if id == "" {
		playerRepo := bf1model.NewPlayerRepository(global.DB)
		// 检查是否已经绑定
		var data *bf1model.Player
		var err error
		if data, err = playerRepo.GetByQID(ctx.Event.UserID); errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("账号未绑定，请使用 .绑定 id 来绑定")
		}
		id = data.DisplayName
	}
	return id, nil
}

// ID2PID 返回pid和id
func ID2PID(qid int64, id string) (string, string, error) {
	var rmu sync.RWMutex
	playerRepo := bf1model.NewPlayerRepository(global.DB)
	var data *bf1model.Player
	var err error
	if id == "" {
		if data, err = playerRepo.GetByQID(qid); errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("账号未绑定，请使用 .绑定 id 来绑定")
		}
		// 若绑定账号时未获取到pid,重新获取并写入数据库
		if data.PersonalID == "" {
			pid, err := api.GetPersonalID(data.DisplayName)
			if err != nil {
				return "", id, errors.New("获取pid失败，请重试")
			}
			rmu.Lock()
			_ = playerRepo.Update(&bf1model.Player{
				Qid:        qid,
				PersonalID: pid,
			})
			rmu.Unlock()
			return pid, id, err
		}
		return data.PersonalID, data.DisplayName, err
	}
	// 检查数据库内是否存在该id
	if data, err = playerRepo.GetByName(id); errors.Is(err, gorm.ErrRecordNotFound) {
		pid, err := api.GetPersonalID(id)
		if err != nil {
			return "", id, errors.New("获取pid失败，请检查id是否有误")
		}
		return pid, id, err
	}
	// 若绑定账号时未获取到pid,重新获取并写入数据库
	if data.PersonalID == "" {
		pid, err := api.GetPersonalID(id)
		if err != nil {
			return "", id, errors.New("获取pid失败，请重试")
		}
		rmu.Lock()
		_ = global.DB.Update(bf1model.Player{
			Qid:        qid,
			PersonalID: pid,
		})
		rmu.Unlock()
		return pid, id, err
	}
	return data.PersonalID, data.DisplayName, err
}

// RequestWeapon 发送武器信息
func RequestWeapon(ctx *zero.Ctx, id, class string) {
	ctx.Send("少女折寿中...")
	pid, id, err := ID2PID(ctx.Event.UserID, id)
	if err != nil {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
		return
	}
	info := &player.Info{Name: id, PersonalID: pid}
	weapon, err := info.GetWeapons(class)
	if err != nil {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
		return
	}
	txt := "id：" + id + "\n"
	wp := ([]player.Weapons)(*weapon)
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
	renderer.Txt2Img(ctx, txt)
}

// GetBF1Recent 获取bf1最近战绩
func GetBF1Recent(id string) (result *player.Recent, err error) {
	u := "https://api.bili22.me/bf1/recent?name=" + id
	data, err := netreq.Request{URL: u}.GetRespBodyBytes()
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	if err != nil {
		return nil, errors.New("ERR: JSON decode failed")
	}
	return result, err
}

// IsValidID 检查id有效性
func IsValidID(id string) (bool, error) {
	vld, err := api.ReturnJSON("https://signin.ea.com/p/ajax/user/checkOriginId?originId="+id, "GET", nil)
	if err != nil {
		return true, errors.New("验证id有效性失败，将继续绑定，请自行检查id是否正确")
	}
	if vld.Get("message").Str != "origin_id_duplicated" {
		return false, nil
	}
	return true, nil
}

// ServerAdminPermission 是否拥有权限
func ServerAdminPermission(ctx *zero.Ctx) bool {
	if zero.AdminPermission(ctx) {
		return true
	}
	groupRepo := bf1model.NewGroupRepository(global.DB)
	adm := groupRepo.IsGroupAdmin(ctx.Event.GroupID, ctx.Event.UserID)
	return adm
}

// ServerOwnerPermission 腐竹权限
func ServerOwnerPermission(ctx *zero.Ctx) bool {
	groupRepo := bf1model.NewGroupRepository(global.DB)
	p := groupRepo.IsGroupOwner(ctx.Event.GroupID, ctx.Event.UserID)
	return p
}
