package main

import (
	_ "beegostudy/conf" //初始化配置信息

	_ "beegostudy/routers" //加载路由

	_ "beegostudy/service/bleve" //初始化索引

	_ "beegostudy/service/cache" //初始化缓存

	"beegostudy/controllers/platform"
	"beegostudy/controllers/platform/data"
	"beegostudy/controllers/template"
	"beegostudy/service/cron" //定时任务准备
	_ "beegostudy/task"       //初始化定时任务
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	go cron.RunCron()
}

func main() {

	//ORM调试模式打开
	beego.BConfig.WebConfig.Session.SessionOn = true //启用Session
	beego.BConfig.Listen.EnableAdmin = true          //启用进程内监控
	beego.BConfig.EnableGzip = true                  //启用gzip压缩

	beego.Router("/login", &platform.LoginController{})
	beego.Router("/logout", &platform.LogoutController{})
	beego.Router("/register", &platform.RegisterController{})

	//数据控制器
	models := map[string]beego.ControllerInterface{
		"system":  &data.SystemController{},
		"site":    &data.SiteController{},
		"catalog": &data.CatalogController{},
		"menu":    &data.MenuController{},
		"user":    &data.UserController{},
		"cron":    &data.CronController{},
		"image":   &data.ImageController{},
		"package": &data.PackageController{},
		"article": &data.ArticleController{},
		"upload":  &data.UploadController{},
		"search":  &data.SearchController{},
		"test":    &data.TestController{},
	}

	for name, controller := range models {
		beego.Router(fmt.Sprintf("/data/%s/:method", name), controller, "*:Get")
	}

	//测试使用
	beego.Router("/data/test/:method", &data.TestController{}, "*:Get")

	//校验用户登录：未登录则重定向到login
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Session("User") == nil {
			if ctx.Input.IsAjax() && ctx.Input.Query("IsSendData") != "" {
				response := make(map[string]interface{})
				response["STATUS"] = 101
				ctx.Output.JSON(response, true, true)
			} else {
				//如果使用dialog方式会出现弹出窗口被定向到了登录页，这里使用js跳转
				//ctx.Redirect(302, "platform.LoginPage")
				ctx.WriteString(platform.LoginPageScript)
			}
			return
		}
	}

	//html格式支持直接预览
	beego.Router("/:path.html", &template.HTMLController{})

	beego.InsertFilter("/platform/*", beego.BeforeRouter, FilterUser)

	if beego.BConfig.RunMode != "dev" {
		beego.InsertFilter("/data/*", beego.BeforeRouter, FilterUser)
	}

	//附件默认目录
	beego.SetStaticPath("/upload", "upload")
	beego.Run()

}
