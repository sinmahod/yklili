package conf

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(60.205.164.3:3306)/beestudy?charset=utf8")
	orm.Debug = true
}
