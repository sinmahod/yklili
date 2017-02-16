package conf

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func RegisterDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(localhost:3306)/beestudy?charset=utf8") //60.205.164.3
	orm.Debug = true
}
