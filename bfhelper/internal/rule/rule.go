// Package rule 命令触发条件
package rule

import (
	_ "embed"
	"encoding/json"
	"io"
	"os"

	fcext "github.com/FloatTech/floatbox/ctxext"
	bf1api "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/bf1/api"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/engine"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/setting"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

//go:embed template.yml
var defaultConfig string

func setupSetting() error {
	setting, err := setting.NewSetting("settings", engine.Engine.DataFolder())
	if err != nil {
		return err
	}
	if err := setting.ReadSection("Account", &global.Account); err != nil {
		return err
	}
	return setting.ReadSection("SakuraKooi", &global.SakuraAPI)
}

func readDictionary() error {
	f, err := os.Open(engine.Engine.DataFolder() + "dic/dic.json")
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &global.Dictionary)
	if err != nil {
		return err
	}
	return nil
}

func generateConfig() {
	logrus.Warnln("[battlefield]未找到配置或者出现错误, 正在重新生成...")
	_ = os.WriteFile(engine.Engine.DataFolder()+"settings.yml", []byte(defaultConfig), 0o644)
	logrus.Warnln("配置已生成! 请修改 settings.yml")
}

// Initialized 需要执行后才能使用插件
func Initialized() zero.Rule {
	return fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		var err error
		dbname := engine.Engine.DataFolder() + "battlefield.db"
		// 初始化数据库
		if err = model.Init(dbname); err != nil {
			ctx.SendChain(message.Text("ERROR: 数据库初始化失败, 请联系机器人管理员重启"))
			return false
		}
		// 读取配置文件
		if err = setupSetting(); err != nil {
			ctx.SendChain(message.Text("ERROR: 读取插件配置失败, 正在重新创建"))
			generateConfig()
			ctx.SendChain(message.Text("INFO: 插件配置已重新创建, 请联系机器人主人修改"))
			return false
		}
		// 建立数据库连接
		global.DB, err = model.Open(dbname)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: 插件数据库连接失败, 请联系机器人管理员重启"))
			return false
		}
		// 刷新Session
		_ = bf1api.Login(global.Account.Username, global.Account.Password)
		// 读字典
		err = readDictionary()
		if err != nil {
			logrus.Errorf("read dictionary: %v", err)
		}
		return true
	})
}

/* 等待数据库重构后也重构
// ServerAdminPermission 是否拥有权限
func ServerAdminPermission() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		if zero.AdminPermission(ctx) {
			return true
		}
			groupRepo := bf1model.NewGroupRepository(global.DB)
			adm := groupRepo.IsGroupAdmin(ctx.Event.GroupID, ctx.Event.UserID)
			return adm
		return false
	}
}

// ServerOwnerPermission 腐竹权限
func ServerOwnerPermission() zero.Rule {
	return func(ctx *zero.Ctx) bool {
			groupRepo := bf1model.NewGroupRepository(global.DB)
			p := groupRepo.IsGroupOwner(ctx.Event.GroupID, ctx.Event.UserID)
			return p
		return true
	}
}
*/
