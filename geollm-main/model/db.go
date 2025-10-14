package model

import (
	"dbdemo/utils"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB
var err error

func InitDb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		utils.DbUser, utils.DbPassWord, utils.DbHost, utils.DbPort, utils.DbName)

	// 添加日志
	logFile, err := os.OpenFile("sql.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	// 将日志输出重定向到文件
	log.SetOutput(logFile)

	Db, err = sqlx.Open(utils.Db, dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}

	// 设置数据库调试
	Db.SetMaxOpenConns(utils.MaxOpenConns)
	Db.SetMaxIdleConns(utils.MaxIdleConns)

	// 测试连接
	if err := Db.Ping(); err != nil {
		fmt.Printf("ping DB failed, err:%v\n", err)
		return
	}

	fmt.Printf("connect DB success\n")
}
