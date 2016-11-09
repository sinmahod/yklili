package platform

type MenusController struct {
	PlatformController
}

func (c *MenusController) Page() {
	c.TplName = "menus.html"
}
