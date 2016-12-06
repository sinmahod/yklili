package main

import (
	"beegostudy/controllers/data"
	"beegostudy/controllers/dml"
	"beegostudy/controllers/platform"
	"beegostudy/util/cron"

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
	eid := this.GetString("EID")
	sid := this.GetString("SID")
	cron.ReadTaskFile()
	if eid != "" {
		cron.StopTask(eid)
	}
	if sid != "" {
		cron.StartTask(sid)
	}

	this.TplName = "test.html"
}

func TestCron() {

	cron.Task("Test", "*/3, *, *, *, *, *", func() {
		fmt.Println("-----------哈哈3------")
	}, "测试定时任务3秒一次")
	cron.Task("Test2", "*/5, *, *, *, *, *", func() {
		fmt.Println("-----------哈哈5------")
	}, "测试定时任务5秒一次")

	cron.StartTasks()

	select {}
}

func main() {
	go TestCron()
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
