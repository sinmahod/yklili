package main

import (
	"beegostudy/models"
	"fmt"
	"github.com/astaxie/beego"
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
	this.Ctx.WriteString("hello world")
}

func main() {

	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	user := new(models.User)
	//o.ReadOrCreate(user, user, ...)
	user.Id = 1
	o.Read(user)

	//user.Phone = "18010182345"
	//user.Name = "gltest"

	fmt.Println(user)

	/*
		    beego.Router("/", &MainController{})
			StaticDir["/static"] = "static"
			beego.Run()
	*/
}
