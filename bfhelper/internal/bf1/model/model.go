// Package bf1model bf1数据库操作
package bf1model

import "github.com/jinzhu/gorm"

// Open 打开数据库
func Open(path string) (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, err
}
