// Package model 数据库操作
package model

import "github.com/jinzhu/gorm"

// Open 打开数据库
func Open(path string) (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, err
}
