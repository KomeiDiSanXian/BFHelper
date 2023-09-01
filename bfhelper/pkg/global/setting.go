package global

import (
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/logger"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/setting"
)

var (
	AccountSetting    *setting.AccountSettingS // AccountSetting 账号设置
	SessionAPISetting *setting.SessionAPISettingS // Session 获取
	BFEACSetting      *setting.BFEACSettingS // BFEAC 设置

	Logger *logger.Logger // Logger 日志
)
