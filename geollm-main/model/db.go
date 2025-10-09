package model

import (
	"dbdemo/utils"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 不要忘了导入数据库驱动
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB
var err error

func InitDb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", utils.DbUser, utils.DbPassWord, utils.DbHost, utils.DbPort, utils.DbName)
	// 也可以使用MustConnect连接不成功就panic
	Db, err = sqlx.Connect(utils.Db, dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	Db.SetMaxOpenConns(utils.MaxOpenConns)
	Db.SetMaxIdleConns(utils.MaxIdleConns)
	fmt.Printf("connect DB success\n")
	return
}
