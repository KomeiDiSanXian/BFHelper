// Package bfhelper 工具
package bfhelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/FloatTech/zbputils/img/text"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	bf1record "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/record"
)

// 简转繁 字典
var twmap map[string]string

// 初始化
func init() {
	// 读字典
	f, err := os.Open(engine.DataFolder() + "dic/dic.json")
	if err != nil {
		logrus.Errorf("open dictionary file failed: %v", err)
		return
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		logrus.Errorf("read dictionary file failed: %v", err)
		return
	}
	err = json.Unmarshal(content, &twmap)
	if err != nil {
		logrus.Errorf("unmarshal dictionary file failed: %v", err)
		return
	}
}

// S2tw 简体转繁体
func S2tw(str string) string {
	result := ""
	for _, v := range str {
		r, ok := twmap[string(v)]
		if ok {
			result += r
		}
	}
	return result
}

// Txt2Img 文字转图片并发送
func Txt2Img(ctx *zero.Ctx, txt string) {
	data, err := text.RenderToBase64(txt, text.FontFile, 400, 20)
	if err != nil {
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("将文字转换成图片时发生错误")))
	}
	if id := ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Image("base64://"+helper.BytesToString(data)))); id.ID() == 0 {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERROR:可能被风控了"))
	}
}

// ReturnBindID 检查是否绑定，返回id
func ReturnBindID(ctx *zero.Ctx, id string) (string, error) {
	if id == "" {
		db, err := bf1model.Open(dbname)
		if err != nil {
			return "", errors.New("打开数据库错误")
		}
		defer db.Close()
		playerRepo := bf1model.NewPlayerRepository(db)
		// 检查是否已经绑定
		var data *bf1model.Player
		if data, err = playerRepo.GetByQID(ctx.Event.UserID); errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("账号未绑定，请使用 .绑定 id 来绑定")
		}
		id = data.DisplayName
	}
	return id, nil
}

// ID2PID 返回pid和id
func ID2PID(qid int64, id string) (string, string, error) {
	db, err := bf1model.Open(dbname)
	if err != nil {
		return "", "", errors.New("打开数据库错误")
	}
	defer db.Close()
	var rmu sync.RWMutex
	playerRepo := bf1model.NewPlayerRepository(db)
	var data *bf1model.Player
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
		_ = db.Update(bf1model.Player{
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
	weapon, err := bf1record.GetWeapons(pid, class)
	if err != nil {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERR：", err))
		return
	}
	txt := "id：" + id + "\n"
	wp := ([]bf1record.Weapons)(*weapon)
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
	Txt2Img(ctx, txt)
}

// GetBF1Recent 获取bf1最近战绩
func GetBF1Recent(id string) (result *bf1record.Recent, err error) {
	u := "https://api.bili22.me/bf1/recent?name=" + id
	data, err := api.ReturnJSON(u, "GET", nil)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(strings.NewReader(data)).Decode(&result)
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
	if gjson.Get(vld, "message").Str != "origin_id_duplicated" {
		return false, nil
	}
	return true, nil
}

// ServerAdminPermission 是否拥有权限
func ServerAdminPermission(ctx *zero.Ctx) bool {
	if zero.AdminPermission(ctx) {
		return true
	}
	db, err := bf1model.Open(dbname)
	if err != nil {
		return false
	}
	groupRepo := bf1model.NewGroupRepository(db)
	adm := groupRepo.IsGroupAdmin(ctx.Event.GroupID, ctx.Event.UserID)
	db.Close()
	return adm
}

// ServerOwnerPermission 腐竹权限
func ServerOwnerPermission(ctx *zero.Ctx) bool {
	db, err := bf1model.Open(dbname)
	if err != nil {
		return false
	}
	groupRepo := bf1model.NewGroupRepository(db)
	p := groupRepo.IsGroupOwner(ctx.Event.GroupID, ctx.Event.UserID)
	db.Close()
	return p
}
