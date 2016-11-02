package controllers

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
)

type PlatformController struct {
	beego.Controller
}

func (c *PlatformController) Get() {
	user := c.GetSession("User").(models.User)
	c.Data["UserName"] = user.GetUserName()
	c.TplName = "platform.html"
}
