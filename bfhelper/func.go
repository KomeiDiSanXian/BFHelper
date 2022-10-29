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
	api "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/api"
	bf1model "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/model"
	bf1record "github.com/KomeiDiSanXian/BFHelper/bfhelper/bf1/record"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"gopkg.in/h2non/gentleman.v2"
	"gorm.io/gorm"
)

// 简转繁 字典
var twmap map[string]string

// 初始化
func init() {
	//读字典
	f, err := os.Open(engine.DataFolder() + "dic/dic.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &twmap)
	if err != nil {
		panic(err)
	}
}

// 查询是否被实锤为外挂
func IsGetBan(id string) bool {
	cli := gentleman.New()
	cli.URL("https://api.gametools.network/bfban/checkban?names=" + id)
	res, err := cli.Request().Send()
	if err != nil {
		return false
	}
	return gjson.Get(res.String(), "names."+strings.ToLower(id)+".hacker").Bool()
}

// 简体转繁体
func S2tw(str string) string {
	result := ""
	for _, v := range str {
		result += twmap[string(v)]
	}
	return result
}

// 文字转图片并发送
func Txt2Img(ctx *zero.Ctx, txt string) {
	data, err := text.RenderToBase64(txt, text.FontFile, 400, 20)
	if err != nil {
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("将文字转换成图片时发生错误")))
	}
	if id := ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Image("base64://"+helper.BytesToString(data)))); id.ID() == 0 {
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERROR:可能被风控了"))
	}
}

// 检查是否绑定，返回id
func ReturnBindID(ctx *zero.Ctx, id string) (string, error) {
	if id == "" {
		gdb, err := bf1model.Open(engine.DataFolder() + "player.db")
		if err != nil {
			return "", errors.New("打开数据库错误")
		}
		db := (*bf1model.PlayerDB)(gdb)
		defer db.Close()
		//检查是否已经绑定
		if data, err := db.FindByQid(ctx.Event.UserID); errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("账号未绑定，请使用 .绑定 id 来绑定")
		} else {
			id = data.DisplayName
		}
	}
	return id, nil
}

// id to pid, 返回pid和id
func ID2PID(qid int64, id string) (string, string, error) {
	gdb, err := bf1model.Open(engine.DataFolder() + "player.db")
	var rmu sync.RWMutex
	if err != nil {
		return "", "", errors.New("打开数据库错误")
	}
	db := (*bf1model.PlayerDB)(gdb)
	defer db.Close()
	if id == "" {
		if data, err := db.FindByQid(qid); errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("账号未绑定，请使用 .绑定 id 来绑定")
		} else {
			//若绑定账号时未获取到pid,重新获取并写入数据库
			if data.PersonalID == "" {
				pid, err := api.GetPersonalID(id)
				if err != nil {
					return "", id, errors.New("获取pid失败，请重试")
				}
				rmu.Lock()
				db.Update(bf1model.Player{
					Qid:        qid,
					PersonalID: pid,
				})
				rmu.Unlock()
				return pid, id, err
			}
			return data.PersonalID, data.DisplayName, err
		}
	} else {
		//检查数据库内是否存在该id
		if data, err := db.FindByName(id); errors.Is(err, gorm.ErrRecordNotFound) {
			pid, err := api.GetPersonalID(id)
			if err != nil {
				return "", id, errors.New("获取pid失败，请检查id是否有误")
			}
			return pid, id, err
		} else {
			//若绑定账号时未获取到pid,重新获取并写入数据库
			if data.PersonalID == "" {
				pid, err := api.GetPersonalID(id)
				if err != nil {
					return "", id, errors.New("获取pid失败，请重试")
				}
				rmu.Lock()
				db.Update(bf1model.Player{
					Qid:        qid,
					PersonalID: pid,
				})
				rmu.Unlock()
				return pid, id, err
			}
			return data.PersonalID, data.DisplayName, err
		}
	}
}

// 发送武器信息
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

// 获取bf1最近战绩
func GetBF1Recent(id string) (result *bf1record.Recent, err error) {
	u := "https://api.bili22.me/bf1/recent?name=" + id
	data, err := api.ReturnJson(u, "GET", nil)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(strings.NewReader(data)).Decode(&result)
	if err != nil {
		return nil, errors.New("ERR: JSON decode failed")
	}
	return result, err
}

// 检查id有效性
func IsValidId(id string) (bool, error) {
	vld, err := api.ReturnJson("https://signin.ea.com/p/ajax/user/checkOriginId?originId="+id, "GET", nil)
	if err != nil {
		return true, errors.New("验证id有效性失败，将继续绑定，请自行检查id是否正确")
	}
	if gjson.Get(vld, "message").Str != "origin_id_duplicated" {
		return false, nil
	}
	return true, nil
}

// 是否拥有权限
func ServerAdminPermission(ctx *zero.Ctx) bool {
	if zero.AdminPermission(ctx) {
		return true
	}
	db, cl, _ := OpenServerDB()
	defer cl()
	adm, _ := db.IsAdmin(ctx.Event.GroupID, ctx.Event.UserID)
	return adm
}

// 腐竹权限
func ServerOwnerPermission(ctx *zero.Ctx) bool {
	db, cl, _ := OpenServerDB()
	defer cl()
	p, _ := db.IsOwner(ctx.Event.GroupID, ctx.Event.UserID)
	return p
}

// 打开服务器数据库，返回数据库及close方法,需要调用close
func OpenServerDB() (*bf1model.ServerDB, func() error, error) {
	sdb, err := bf1model.Open(engine.DataFolder() + "server.db")
	if err != nil {
		sql, _ := sdb.DB()
		sql.Close()
		return nil, nil, errors.New("数据库错误")
	}
	db := (*bf1model.ServerDB)(sdb)
	return db, db.Close, nil
}
