package handler

import (
	"context"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/service"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// CreateGroupHandler 创建服务器群组处理函数
func CreateGroupHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.CreateGroup()
	})
}

// DeleteGroupHandler 所在群删除服务器群组处理函数
func DeleteGroupHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.DeleteGroup()
	})
}

// ChangeOwnerHandler 更换服务器群组所有人处理函数
func ChangeOwnerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.ChangeOwner()
	})
}

// AddServerHandler 添加服务器处理函数
func AddServerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.AddServer()
	})
}

// AddServerAdminHandler 添加服务器管理员
func AddServerAdminHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.AddServerAdmin()
	})
}

// SetServerAliasHandler 设置服务器别名
func SetServerAliasHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.SetServerAlias()
	})
}

// DeleteServerHandler 删除服务器处理函数
func DeleteServerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.DeleteServer()
	})
}

// DeleteAdminHandler 删除群组服务器管理员
func DeleteAdminHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.DeleteAdmin()
	})
}

// KickPlayerHandler 踢出玩家处理
func KickPlayerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.KickPlayer(ctx)
	})
}

// BanPlayerHandler 单服务器封禁玩家处理
func BanPlayerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.BanPlayer(ctx)
	})
}

// BanPlayerAtAllServerHandler 已绑定服务器封禁玩家
func BanPlayerAtAllServerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.BanPlayerAtAllServer(ctx)
	})
}

// UnbanPlayerHandler 单服务器解封玩家处理
func UnbanPlayerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.UnbanPlayer(ctx)
	})
}

// UnbanPlayerAtAllServerHandler 已绑定服务器解封玩家
func UnbanPlayerAtAllServerHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(ctx context.Context, svc *service.Service) error {
		return svc.UnbanPlayerAtAllServer(ctx)
	})
}

// ChangeMapHandler 切换地图处理
func ChangeMapHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.ChangeMap()
	})
}

// ReadMapsHandler 查看地图池处理函数
func ReadMapsHandler(zctx *zero.Ctx) {
	GenericHandler(zctx, func(_ context.Context, svc *service.Service) error {
		return svc.GetMap()
	})
}
