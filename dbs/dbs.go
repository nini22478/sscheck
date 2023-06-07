package dbs

import (
	"check_vpn/dbs/query"
	"check_vpn/mylog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const MySQLDSN = "username:passwd@(host:port)/tablename?charset=utf8mb4&parseTime=True&loc=Local"

var adb *gorm.DB

func GetDb() *query.Query {
	if adb == nil {
		var err error
		adb, err = gorm.Open(mysql.Open(MySQLDSN))
		mylog.Der(err)
		query.SetDefault(adb)

	}
	return query.Q
}
