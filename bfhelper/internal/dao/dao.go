// Package dao 数据访问对象
package dao

import "github.com/jinzhu/gorm"

// Dao 含有数据库对象
type Dao struct {
	engine *gorm.DB
}

// New 新建一个数据库对象
func New(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}
