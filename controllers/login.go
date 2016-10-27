package controllers

import (
	"fmt"

	"beegostudy/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	u := c.GetString("u")
	//p := c.GetString("p")

	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	user := new(models.User)
	user.UserName = u
	o.Read(user)
	fmt.Println(user)

	c.Data["UserName"] = user.String()
	c.TplName = "login.html"
}
