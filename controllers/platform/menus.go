package platform

type MenusController struct {
	PlatformController
}

func (c *MenusController) Get() {
	c.TplName = "menus.html"
}
