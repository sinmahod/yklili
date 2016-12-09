package main

import (
	"beegostudy/controllers/data"
	"beegostudy/controllers/dml"
	"beegostudy/controllers/platform"
	"beegostudy/service/cron"
	_ "beegostudy/task"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(60.205.164.3:3306)/beestudy?charset=utf8")

	go cron.RunCron()
}

var Crontab = cron.New()

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
		"user": &data.UserController{},
		"cron": &data.CronController{},
	}

	for name, controller := range models {
		beego.Router(fmt.Sprintf("/data/%s/:method", name), controller, "*:Get")
	}

	beego.Router("/platform/test", &platform.TestController{})

	//校验用户登录：未登录则重定向到login
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Session("User") == nil {
			//如果使用dialog方式会出现弹出窗口被定向到了登录页，这里使用js跳转
			//ctx.Redirect(302, "platform.LoginPage")
			ctx.WriteString(platform.LoginPageScript)
		}
	}

	//dml格式支持直接预览
	beego.AddTemplateExt("dml")
	beego.Router("/:path.dml", &dml.DMLController{})

	beego.InsertFilter("/platform/*", beego.BeforeRouter, FilterUser)
	beego.Run()

}
