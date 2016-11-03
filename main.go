package main

import (
	"beegostudy/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(60.205.164.3:3306)/beestudy?charset=utf8")
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["UserName"] = "HHHHH"
	this.TplName = "test.html"
}

func main() {
	orm.Debug = true                                 //ORM调试模式打开
	beego.BConfig.WebConfig.Session.SessionOn = true //启用Session

	beego.Router("/", &MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})
	beego.Router("/platform", &controllers.PlatformController{})

	//校验用户登录：未登录则重定向到login
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Session("User") == nil {
			ctx.Redirect(302, "/login")
		}
	}

	beego.InsertFilter("/platform/*", beego.BeforeRouter, FilterUser)
	beego.Run()

	/**************自动建表***********/
	// //数据库别名
	// name := "default"

	// // drop table 后再建表
	// force := true

	// // 打印执行过程
	// verbose := true

	// // 遇到错误立即返回
	// err := orm.RunSyncdb(name, force, verbose)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
