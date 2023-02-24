package controller

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {

	username := "root"     //账号
	password := "mm123456" //密码
	host := "127.0.0.1"    //数据库地址，可以是Ip或者域名
	port := 10112          //数据库端口
	Dbname := "douyin"     //数据库名
	timeout := "10s"       //连接超时，10秒

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("connect to mysql failed, error=" + err.Error())
	} else {
		DB = db
		DB.AutoMigrate(&UserDto{})  // 如果表没建立，建立该表
		DB.AutoMigrate(&VideoDto{}) // 如果表没建立，建立该表
		DB.AutoMigrate(&Favorite{}) // 如果表没建立，建立该表
		DB.AutoMigrate(&Comment{})  // 如果表没建立，建立该表
	}
	return db, err
}
