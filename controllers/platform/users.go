package platform

type UsersController struct {
	PlatformController
}

func (c *UsersController) Page() {
	c.TplName = "platform/users.html"
}
