// Package rule 命令触发条件
package rule

import (
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/dao"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/model"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/logger"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/setting"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:embed template.yml
var defaultConfig string

//go:embed dic.json
var traditionalChinese string

func init() {
	dbname := global.Engine.DataFolder() + "battlefield.db"
	err := model.Init(dbname)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	global.DB, err = model.Open(dbname)
	if err != nil {
		logrus.Fatalf("Failed to open database: %v", err)
	}
	if err := setupSetting(); err != nil {
		generateConfig()
		os.Exit(1)
	}
	// 读字典
	err = readDictionary()
	if err != nil {
		logrus.Errorf("read dictionary: %v", err)
	}
}

func setupSetting() error {
	settings, err := setting.NewSetting("settings", global.Engine.DataFolder())
	if err != nil {
		return err
	}
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
