package platform

type UsersController struct {
	PlatformController
}

func (c *UsersController) Page() {
	c.TplName = "users.html"
}
