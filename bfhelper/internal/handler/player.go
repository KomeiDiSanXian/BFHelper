// Package handler 事件处理函数
package handler

import (
	"context"
	"reflect"
	"runtime"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/errcode"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/service"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"go.opentelemetry.io/otel/codes"
)

// GenericHandler 通用处理封装
func GenericHandler(zctx *zero.Ctx, serviceMethod func(context.Context, *service.Service) error) {
	svc := service.New(zctx)
	funcName := runtime.FuncForPC(reflect.ValueOf(serviceMethod).Pointer()).Name()
	nCtx, span := global.Tracer.Start(context.Background(), "Handler")
	err := serviceMethod(nCtx, svc)
	defer span.End()

	if errors.Is(err, errcode.Success) || errors.Is(err, errcode.Canceled) {
		span.SetStatus(codes.Ok, "")
		return
	}
	logrus.Errorf("%s error: %v", funcName, err)
	svc.Log().Error(err)
	span.SetStatus(codes.Error, "received error")
	span.RecordError(err)
}

// BindAccountHandler 绑定账号处理函数
func BindAccountHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.BindAccount(ctx)
	})
}

// PlayerRecentHandler 最近战绩查询处理函数
func PlayerRecentHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.GetPlayerRecent(ctx)
	})
}

// PlayerStatsHandler 玩家战绩查询处理函数
func PlayerStatsHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.GetPlayerStats(ctx)
	})
}

// PlayerWeaponHandler 玩家武器查询处理函数
func PlayerWeaponHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.GetPlayerWeapon(ctx)
	})
}

// PlayerVehicleHandler 玩家载具查询处理函数
func PlayerVehicleHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.GetPlayerVehicle(ctx)
	})
}

// BF1ExchangeHandler 获取战地一本期交换处理函数
func BF1ExchangeHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.GetBF1Exchange(ctx)
	})
}

// BF1OpreationPackHandler 获取战地一本期行动包处理函数
func BF1OpreationPackHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.GetBF1OpreationPack(ctx)
	})
}

// PlayerBanInfoHandler 获取玩家联ban信息
func PlayerBanInfoHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.GetPlayerBanInfo(ctx)
	})
}
