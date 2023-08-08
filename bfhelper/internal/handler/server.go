package handler

import (
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/service"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// CreateGroupHandler 创建服务器群组处理函数
func CreateGroupHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.CreateGroup()
	})
}

// DeleteGroupHandler 所在群删除服务器群组处理函数
func DeleteGroupHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.DeleteGroup()
	})
}

// ChangeOwnerHandler 更换服务器群组所有人处理函数
func ChangeOwnerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.ChangeOwner()
	})
}

// AddServerHandler 添加服务器处理函数
func AddServerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.AddServer()
	})
}

// AddServerAdminHandler 添加服务器管理员
func AddServerAdminHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.AddServerAdmin()
	})
}

// SetServerAliasHandler 设置服务器别名
func SetServerAliasHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.SetServerAlias()
	})
}

// DeleteServerHandler 删除服务器处理函数
func DeleteServerHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.DeleteServer()
	})
}

// DeleteAdminHandler 删除群组服务器管理员
func DeleteAdminHandler() zero.Handler {
	return ErrorHandlerWrapper(func(svc *service.Service) error {
		return svc.DeleteAdmin()
	})
}
