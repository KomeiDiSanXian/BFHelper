package bfhelper

import (
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/engine"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/handler"
	"github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/rule"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	// 群聊绑定多服务器
	//
	// 每个群聊可以绑定多个服务器, 每个服务器也可以绑定多个群聊
	engine.Engine.OnPrefix(".创建服务器群组", zero.OnlyGroup, zero.OwnerPermission, rule.Initialized()).SetBlock(true).Handle(handler.CreateGroupHandler())
	engine.Engine.OnFullMatch(".删除服务器群组", zero.OnlyGroup, zero.OwnerPermission, rule.Initialized()).SetBlock(true).Handle(handler.DeleteGroupHandler())
	engine.Engine.OnPrefix(".更换服主", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.ChangeOwnerHandler())
	engine.Engine.OnPrefix(".绑定服务器", zero.OnlyGroup, rule.Initialized(), zero.SuperUserPermission).SetBlock(true).Handle(handler.AddServerHandler())
	engine.Engine.OnPrefix(".添加管理", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.AddServerAdminHandler())
	engine.Engine.OnPrefix(".设置别名", zero.OnlyGroup, rule.Initialized(), rule.ServerAdminPermission()).SetBlock(true).Handle(handler.SetServerAliasHandler())
	engine.Engine.OnPrefix(".解绑服务器", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.DeleteServerHandler())
	engine.Engine.OnPrefix(".删除管理", zero.OnlyGroup, rule.Initialized(), rule.ServerOwnerPermission()).SetBlock(true).Handle(handler.DeleteAdminHandler())
}
