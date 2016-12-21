package main

import (
	"beegostudy/controllers/platform"
	"beegostudy/controllers/platform/data"
	"beegostudy/controllers/template"
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
	beego.BConfig.Listen.EnableAdmin = true          //启用进程内监控

	beego.Router("/", &MainController{})
	beego.Router("/login", &platform.LoginController{})
	beego.Router("/register", &platform.RegisterController{})

	//数据控制器
	models := map[string]beego.ControllerInterface{
		"menu":  &data.MenuController{},
		"user":  &data.UserController{},
		"cron":  &data.CronController{},
		"image": &data.ImageController{},
		"test":  &data.TestController{},
	}

	for name, controller := range models {
		beego.Router(fmt.Sprintf("/data/%s/:method", name), controller, "*:Get")
	}

	//测试使用
	beego.Router("/data/test/:method", &data.TestController{}, "*:Get")

	//校验用户登录：未登录则重定向到login
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Session("User") == nil {
			//如果使用dialog方式会出现弹出窗口被定向到了登录页，这里使用js跳转
			//ctx.Redirect(302, "platform.LoginPage")
			ctx.WriteString(platform.LoginPageScript)
		}
	}

	//html格式支持直接预览
	beego.Router("/:path.html", &template.HTMLController{})

	beego.InsertFilter("/platform/*", beego.BeforeRouter, FilterUser)

	//附件默认目录
	beego.SetStaticPath("/upload", "upload")
	beego.Run()

}
