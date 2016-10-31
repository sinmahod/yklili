package controllers

import (
	"fmt"

	//"beegostudy/models"

	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	u := c.GetString("username")
	p := c.GetString("password")
	fmt.Println(u, p)
	if u != "" {
		c.TplName = "test.html"
	} else {
		c.TplName = "login.html"
	}
}

func (c *LoginController) Post() {

	//c.Ctx.Redirect(200, "test.html")
}
