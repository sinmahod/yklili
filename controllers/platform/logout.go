package platform

import (
	"github.com/astaxie/beego"
)

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Get() {
	c.DelSession("User")
	c.TplName = "login.html"
}
