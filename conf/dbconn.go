package conf

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var (
	DBType     string
	DBLocal    string
	DBPort     string
	DBName     string
	DBUserName string
	DBPassword string
)

func RegisterDB() {
	DBType = strings.ToLower(beego.AppConfig.String("dbtype"))
	DBLocal = beego.AppConfig.String("dblocal")
	DBPort = beego.AppConfig.String("dbport")
	DBName = beego.AppConfig.String("dbname")
	DBUserName = beego.AppConfig.String("dbusername")
	DBPassword = beego.AppConfig.String("dbpassword")

	orm.RegisterDriver(DBType, orm.DRMySQL)
	if DBType == "mysql" {
		orm.RegisterDataBase("default", DBType, DBUserName+":"+DBPassword+"@tcp("+DBLocal+":"+DBPort+")/"+DBName+"?charset=utf8") //60.205.164.3
	}
	mode := beego.AppConfig.String("runmode")

	if mode == "dev" {
		orm.Debug = true
	}

}
