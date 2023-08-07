// Package handler 事件处理函数
package handler

import (
	"reflect"
	"runtime"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/service"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// ErrorHandlerWrapper 错误处理封装
func ErrorHandlerWrapper(serviceMethod func(*service.Service) error) zero.Handler {
	return func(ctx *zero.Ctx) {
		svc := service.New(ctx)
		err := serviceMethod(svc)
		// TODO: 写入日志
		if err != nil {
			logrus.Errorf("%s error: %v", runtime.FuncForPC(reflect.ValueOf(serviceMethod).Pointer()).Name(), err)
		}
	}
}

// BindAccountHandler 绑定账号处理函数
func BindAccountHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.BindAccount()
	})
}

// PlayerRecentHandler 最近战绩查询处理函数
func PlayerRecentHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.GetPlayerRecent()
	})
}

// PlayerStatsHandler 玩家战绩查询处理函数
func PlayerStatsHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.GetPlayerStats()
	})
}

// PlayerWeaponHandler 玩家武器查询处理函数
func PlayerWeaponHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.GetPlayerWeapon()
	})
}

// PlayerVehicleHandler 玩家载具查询处理函数
func PlayerVehicleHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.GetPlayerVehicle()
	})
}

// BF1ExchangeHandler 获取战地一本期交换处理函数
func BF1ExchangeHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.GetBF1Exchange()
	})
}

// BF1OpreationPackHandler 获取战地一本期行动包处理函数
func BF1OpreationPackHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.GetBF1OpreationPack()
	})
}

// PlayerBanInfoHandler 获取玩家联ban信息
func PlayerBanInfoHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.GetPlayerBanInfo()
	})
}
