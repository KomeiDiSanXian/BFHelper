// Package model 数据库操作
package model

import (
	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite 启用dialect，集成测试及relase需要注释掉
)

// Init 数据库初始化
func Init(path string) error {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return err
	}

	// Migrate the schema
	err = db.AutoMigrate(&Player{}, &Group{}, &Server{}, &Admin{}).Error
	if err != nil {
		return err
	}
	db.SingularTable(true)
	return db.Close()
}
