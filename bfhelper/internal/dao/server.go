// Package dao 服务器群组dao层操作
package dao

import "github.com/KomeiDiSanXian/BFHelper/bfhelper/internal/model"

// CreateGroup 创建新的服务器群组
func (d *Dao) CreateGroup(groupID, owner int64) error {
	grp := model.NewGroup(groupID)
	grp.Owner = owner
	return grp.Create(d.engine)
}

// DeleteGroup 删除服务器群组
func (d *Dao) DeleteGroup(groupID int64) error {
	return model.NewGroup(groupID).DeleteByID(d.engine)
}

// UpdateOwner 更新服务器服主
func (d *Dao) UpdateOwner(groupID, owner int64) error {
	grp := model.NewGroup(groupID)
	grp.Owner = owner
	return grp.Update(d.engine)
}

// AddGroupServer 添加服务器
func (d *Dao) AddGroupServer(groupID int64, gameID, serverID, pgid, serverName string) error {
	srv := model.NewServer(gameID)
	srv.ServerID, srv.PGID, srv.ServerName = serverID, pgid, serverName
	grp := model.NewGroup(groupID)
	grp.Servers = append(grp.Servers, *srv)
	return grp.Update(d.engine)
}

// AddGroupAdmin 添加服务器管理员
func (d *Dao) AddGroupAdmin(groupID int64, adminQQ ...int64) error {
	grp := model.NewGroup(groupID)
	for _, qq := range adminQQ {
		grp.Admins = append(grp.Admins, *model.NewAdmin(qq))
	}
	return grp.Update(d.engine)
}

// ServerAlias 服务器别名修改
func (d *Dao) ServerAlias(gameID, alias string) error {
	srv := model.NewServer(gameID)
	srv.NameInGroup = alias
	return srv.Update(d.engine)
}

// RemoveGroupServer 移除群组服务器
func (d *Dao) RemoveGroupServer(groupID int64, gameID string) error {
	return model.NewGroup(groupID).DeleteServer(d.engine, gameID)
}

// RemoveGroupAdmin 移除服务器管理员
func (d *Dao) RemoveGroupAdmin(groupID, adminQQ int64) error {
	return model.NewGroup(groupID).DeleteAdmin(d.engine, adminQQ)
}

// IsServerAdmin 判断是否为服务器管理
func (d *Dao) IsServerAdmin(groupID, qq int64) bool {
	return model.NewGroup(groupID).IsAdmin(d.engine, qq)
}

// IsOwner 判断是不是服务器服主
func (d *Dao) IsOwner(groupID, qq int64) bool {
	return model.NewGroup(groupID).IsOwner(d.engine, qq)
}

// GetGroup 获取群组
func (d *Dao) GetGroup(groupID int64) (*model.Group, error) {
	return model.NewGroup(groupID).GetByID(d.engine)
}

// GetServerByAlias 通过别名获取服务器信息
func (d *Dao) GetServerByAlias(groupID int64, alias string) (*model.Server, error) {
	return model.NewGroup(groupID).GetByAlias(d.engine, alias)
}
