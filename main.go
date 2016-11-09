package main

import (
	"beegostudy/controllers/data"
	"beegostudy/controllers/platform"

	"fmt"
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

	//页面控制器
	pages := map[string]beego.ControllerInterface{
		"users": &platform.UsersController{},
		"menus": &platform.MenusController{},
	}
	for name, controller := range pages {
		beego.Router(fmt.Sprintf("/platform/%s", name), controller, "get:Page")
	}

	//数据控制器
	models := map[string]beego.ControllerInterface{
		"menu": &data.MenuController{},
	}
	for name, controller := range models {
		beego.Router(fmt.Sprintf("/data/%s/:method", name), controller)
	}

	beego.Router("/platform/test", &platform.TestController{})

	//校验用户登录：未登录则重定向到login
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Session("User") == nil {
			ctx.Redirect(302, "/login")
		}
	}

	beego.InsertFilter("/platform/*", beego.BeforeRouter, FilterUser)
	beego.Run()
}
