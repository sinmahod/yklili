package platform

type MenusController struct {
	PlatformController
}

func (c *MenusController) Page() {
	c.TplName = "platform/menus.html"
}
