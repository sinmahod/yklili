package main

import (
	"beegostudy/controllers/platform"
	"beegostudy/controllers/platform/data"

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
	beego.Router("/login", &platform.LoginController{})
	beego.Router("/register", &platform.RegisterController{})
	beego.Router("/platform/users", &platform.UsersController{})
	beego.Router("/platform/menus", &platform.MenusController{})
	beego.Router("/platform/test", &platform.TestController{})

	beego.Router("/data/menulist", &data.MenuListController{})

	//校验用户登录：未登录则重定向到login
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Session("User") == nil {
			ctx.Redirect(302, "/login")
		}
	}

	beego.InsertFilter("/platform/*", beego.BeforeRouter, FilterUser)
	beego.Run()
}
