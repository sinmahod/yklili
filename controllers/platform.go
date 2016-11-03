package controllers

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
)

type PlatformController struct {
	beego.Controller
}

func (c *PlatformController) Get() {

	c.Layout = "platform.html"

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Include"] = "public/include.tpl"
	c.LayoutSections["Script"] = "public/script.tpl"

	user := c.GetSession("User").(models.User)
	c.Data["UserName"] = user.GetUserName()
	c.TplName = "content.html"
}
