// Package bf1model bf1数据库操作
package bf1model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	return db.Close()
}
