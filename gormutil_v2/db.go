package gormutil

import (
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var defaultDB *gorm.DB

// DB 为了 UnitTest 可以 mock
var DB func() *gorm.DB = globalDB

func globalDB() *gorm.DB {
	return defaultDB
}

// 通过 dsn, conf 创建 db 连接
func ConnectWithDSN(dsn string, conf Conf) (*gorm.DB, error) {
	newLogger := logger.Default.LogMode(logger.Warn)
	if conf.LogMode {
		newLogger.LogMode(logger.Info)
	}
	var (
		db  *gorm.DB
		err error
	)

	if strings.HasPrefix(dsn, "postgres://") {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	} else if strings.HasPrefix(dsn, "mysql://") {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	} else {
		return nil, gorm.ErrUnsupportedDriver
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	if conf.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	}

	if conf.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	}

	if conf.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second)
	}

	return db, nil
}

// 通过 Conf 创建全局 db
func ConnectGlobalDB(conf Conf) error {
	db, err := ConnectWithDSN(conf.DSN, conf)

	if err != nil {
		return err
	}

	defaultDB = db

	return nil
}

// 使用 conf.C 创建全局 db 的便捷方式
func Connect() error {
	return ConnectGlobalDB(C)
}
