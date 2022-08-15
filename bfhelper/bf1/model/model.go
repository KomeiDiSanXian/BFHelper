package bf1model

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

//如果数据库不存在就创建数据库
func InitDB(dbpath string, tables ...interface{}) (*gorm.DB, error) {
	if _, err := os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	}
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(tables...)
	return db, nil
}

//打开数据库
func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	return db, err
}
