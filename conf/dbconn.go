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
	beego.Info("正在初始化数据库连接...")
	DBType = strings.ToLower(beego.AppConfig.String("DB.Type"))
	DBLocal = beego.AppConfig.String("DB.IP")
	DBPort = beego.AppConfig.String("DB.Port")
	DBName = beego.AppConfig.String("DB.Name")
	DBUserName = beego.AppConfig.String("DB.UserName")
	DBPassword = beego.AppConfig.String("DB.Password")

	orm.RegisterDriver(DBType, orm.DRMySQL)
	if DBType == "mysql" {
		err := orm.RegisterDataBase("default", DBType, DBUserName+":"+DBPassword+"@tcp("+DBLocal+":"+DBPort+")/"+DBName+"?charset=utf8") //60.205.164.3
		if err != nil {
			beego.Error("数据库连接失败：", err)
			return
		}
	}

	beego.Info("数据库连接成功")

	mode := beego.AppConfig.String("runmode")

	if mode == "dev" {
		beego.Info("当前运行环境为开发环境，将在控制台打印SQL语句")
		orm.Debug = true
	}

}
