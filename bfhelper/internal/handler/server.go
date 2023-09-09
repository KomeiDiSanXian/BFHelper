package handler

import (
	"context"

	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/service"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// CreateGroupHandler 创建服务器群组处理函数
func CreateGroupHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.CreateGroup()
	})
}

// DeleteGroupHandler 所在群删除服务器群组处理函数
func DeleteGroupHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.DeleteGroup()
	})
}

// ChangeOwnerHandler 更换服务器群组所有人处理函数
func ChangeOwnerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.ChangeOwner()
	})
}

// AddServerHandler 添加服务器处理函数
func AddServerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.AddServer()
	})
}

// AddServerAdminHandler 添加服务器管理员
func AddServerAdminHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.AddServerAdmin()
	})
}

// SetServerAliasHandler 设置服务器别名
func SetServerAliasHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.SetServerAlias()
	})
}

// DeleteServerHandler 删除服务器处理函数
func DeleteServerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.DeleteServer()
	})
}

// DeleteAdminHandler 删除群组服务器管理员
func DeleteAdminHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.DeleteAdmin()
	})
}

// KickPlayerHandler 踢出玩家处理
func KickPlayerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.KickPlayer()
	})
}

// BanPlayerHandler 单服务器封禁玩家处理
func BanPlayerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.BanPlayer()
	})
}

// BanPlayerAtAllServerHandler 已绑定服务器封禁玩家
func BanPlayerAtAllServerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.BanPlayerAtAllServer()
	})
}

// UnbanPlayerHandler 单服务器解封玩家处理
func UnbanPlayerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.UnbanPlayer()
	})
}

// UnbanPlayerAtAllServerHandler 已绑定服务器解封玩家
func UnbanPlayerAtAllServerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.UnbanPlayerAtAllServer()
	})
}

// ChangeMapHandler 切换地图处理
func ChangeMapHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.ChangeMap()
	})
}

// ReadMapsHandler 查看地图池处理函数
func ReadMapsHandler() zero.Handler {
	return ErrorHandlerWrapper(func(ctx context.Context, svc *service.Service) error {
		return svc.GetMap()
	})
}
