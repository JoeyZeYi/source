//go:build linux
// +build linux

package data

import (
	"fmt"
	"github.com/JoeyZeYi/source/log"
	"github.com/JoeyZeYi/source/log/zap"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewGormDB(addr, userName, pwd, dbName string, openConn, idleConn int, gormLogger gormlogger.Interface) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, pwd, addr, dbName)
	directory := gormMysql.Open(dns)

	//gormLogger := log.NewGormLogger(log.GormLoggerSlowThreshold(time.Second), log.GormLoggerLevel(gormlogger.Info))
	db, err := gorm.Open(directory, &gorm.Config{
		Logger: gormLogger,
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
	sqlDb.SetMaxOpenConns(openConn)
	sqlDb.SetMaxIdleConns(idleConn)
	return db, nil
}
