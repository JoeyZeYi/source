//go:build windows
// +build windows

package data

import (
	"fmt"
	"github.com/JoeyZeYi/source/log"
	"github.com/JoeyZeYi/source/log/zap"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(addr, userName, pwd, dbName string, openConn, idleConn int, l logger.Interface) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", userName, pwd, addr, dbName)
	directory := gormMysql.Open(dns)

	//logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(directory, &gorm.Config{
		Logger: l,
	})
	if err != nil {
		log.Error("getGormDB", zap.Error(err))
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		log.Error("getGormDB", zap.Error(err))
		return nil, err
	}
	//默认都是200
	sqlDb.SetMaxOpenConns(openConn)
	sqlDb.SetMaxIdleConns(idleConn)
	return db, nil
}
