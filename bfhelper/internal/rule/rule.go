// Package rule 命令触发条件
package rule

import (
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/dao"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/logger"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/setting"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:embed template.yml
var defaultConfig string

//go:embed dic.json
var traditionalChinese string

var settings *setting.Setting

func init() {
	dbname := global.Engine.DataFolder() + "battlefield.db"
	_ = model.Init(dbname)
	global.DB, _ = model.Open(dbname)
	settings, _ = setting.NewSetting("settings", global.Engine.DataFolder())
}

func setupSetting() error {
	if err := settings.ReadSection("Account", &global.AccountSetting); err != nil {
		return err
	}
	if err := settings.ReadSection("BFEAC", &global.BFEACSetting); err != nil {
		return err
	}
	if err := settings.ReadSection("Trace", &global.TraceSetting); err != nil {
		return err
	}
	setupLogger()
	return settings.ReadSection("SakuraKooi", &global.SessionAPISetting)
}

func setupLogger() {
	fileName := global.Engine.DataFolder() + "/log/bfhelper.log"
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   1024, // 最大为1G
		MaxAge:    10,   // 日志保存10天
		LocalTime: true,
	}, "", log.LstdFlags)
}

func readDictionary() error {
	content, err := io.ReadAll(strings.NewReader(traditionalChinese))
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
	_ = os.WriteFile(global.Engine.DataFolder()+"settings.yml", []byte(defaultConfig), 0o644)
	logrus.Warnln("配置已生成! 请修改 settings.yml")
}

// Initialized 需要执行后才能使用插件
func Initialized(ctx *zero.Ctx) bool {
	rule := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		var err error
		// 读取配置文件
		if err = setupSetting(); err != nil {
			ctx.SendChain(message.Text("ERROR: 读取插件配置失败, 正在重新创建"))
			generateConfig()
			ctx.SendChain(message.Text("INFO: 插件配置已重新创建, 请联系机器人主人修改"))
			return false
		}
		// 读字典
		err = readDictionary()
		if err != nil {
			logrus.Errorf("read dictionary: %v", err)
		}
		return true
	})
	return rule(ctx)
}

// ServerAdminPermission 是否拥有权限
func ServerAdminPermission(ctx *zero.Ctx) bool {
	if zero.AdminPermission(ctx) {
		return true
	}
	return dao.New(global.DB).IsServerAdmin(ctx.Event.GroupID, ctx.Event.UserID)
}

// ServerOwnerPermission 腐竹权限
func ServerOwnerPermission(ctx *zero.Ctx) bool {
	if zero.OwnerPermission(ctx) {
		return true
	}
	return dao.New(global.DB).IsOwner(ctx.Event.GroupID, ctx.Event.UserID)
}
