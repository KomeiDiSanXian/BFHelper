package bf1model

import (
	"os"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 读写锁
var rmu sync.RWMutex

// 如果数据库不存在就创建数据库
func InitDB(dbpath string, tables ...interface{}) error {
	if _, err := os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return err
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
	sqlDB, _ := db.DB()
	return sqlDB.Close()
}

// 打开数据库
func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	return db, err
}
