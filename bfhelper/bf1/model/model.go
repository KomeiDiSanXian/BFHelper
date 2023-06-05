package bf1model

import "github.com/jinzhu/gorm"

// Open 打开数据库
func Open(path string) (db *gorm.DB, closedb func() error, err error) {
	db, err = gorm.Open("sqlite3", path)
	if err != nil {
		return nil, nil, err
	}
	return db, db.Close, err
}
