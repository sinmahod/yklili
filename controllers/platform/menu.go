package platform

import (
//"github.com/astaxie/beego"
)

type MenusController struct {
	PlatformController
}

func (c *MenusController) Get() {
	c.TplName = "menus.html"
}
