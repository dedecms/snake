package snake

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database ...
type Database interface {
	Get(str string) (db *gorm.DB, err error) // 返回数据库链接
}

// DB Database 初始化 ...
func DB(conf string) snakeDatabase {
	var db snakeDatabase
	c := Text(conf).Split("@")
	u := Text(c[0]).Split(":")
	db.User = u[0]
	db.Pass = u[1]
	h := Text(c[1]).Split(":")
	db.Host = h[0]
	db.Port = h[1]
	return db
}

type snakeDatabase struct {
	DBname []string
	User   string
	Pass   string
	Host   string
	Port   string
}

// Add 添加新的数据库 ...
func (sk *snakeDatabase) Get(str string) (db *gorm.DB, err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 5000 * time.Second, // Slow SQL threshold
		},
	)
	return gorm.Open(mysql.New(mysql.Config{
		DSN: Text(sk.User).
			Add(":").
			Add(sk.Pass).
			Add("@tcp(").
			Add(sk.Host).
			Add(":").
			Add(sk.Port).
			Add(")/").
			Add(str).
			Add("?charset=utf8&parseTime=True&loc=Local").
			Get(), // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: newLogger,
	})
}
