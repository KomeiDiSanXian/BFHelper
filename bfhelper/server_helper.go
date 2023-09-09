package bfhelper

import (
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/handler"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/rule"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/pkg/global"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	// 群聊绑定多服务器
	//
	// 每个群聊可以绑定多个服务器, 每个服务器也可以绑定多个群聊
	global.Engine.OnPrefix(".创建服务器群组", zero.OnlyGroup, zero.OwnerPermission, rule.Initialized()).SetBlock(true).Handle(handler.CreateGroupHandler())
	global.Engine.OnFullMatch(".删除服务器群组", zero.OnlyGroup, zero.OwnerPermission, rule.Initialized()).SetBlock(true).Handle(handler.DeleteGroupHandler())
	global.Engine.OnPrefix(".更换服主", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.ChangeOwnerHandler())
	global.Engine.OnPrefix(".绑定服务器", zero.OnlyGroup, rule.Initialized(), zero.SuperUserPermission).SetBlock(true).Handle(handler.AddServerHandler())
	global.Engine.OnPrefix(".添加管理", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.AddServerAdminHandler())
	global.Engine.OnPrefix(".设置别名", zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.SetServerAliasHandler())
	global.Engine.OnPrefix(".解绑服务器", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.DeleteServerHandler())
	global.Engine.OnPrefix(".删除管理", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.DeleteAdminHandler())

	global.Engine.OnPrefixGroup([]string{".踢出", ".kick", ".k"}, zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.KickPlayerHandler())
	global.Engine.OnPrefixGroup([]string{".封禁", ".b"}, zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.BanPlayerHandler())
	global.Engine.OnPrefixGroup([]string{".解封", ".ub"}, zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.UnbanPlayerHandler())
	global.Engine.OnPrefixGroup([]string{".全封", ".bana"}, zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.BanPlayerAtAllServerHandler())
	global.Engine.OnPrefixGroup([]string{".全解", ".ubana"}, zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.UnbanPlayerAtAllServerHandler())
	global.Engine.OnPrefixGroup([]string{".切图", ".cm"}, zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.ChangeMapHandler())
	global.Engine.OnPrefixGroup([]string{".查图", ".qm"}, zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.ReadMapsHandler())
}
