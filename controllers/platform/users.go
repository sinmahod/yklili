package platform

type UsersController struct {
	PlatformController
}

func (c *UsersController) Get() {
	c.TplName = "users.html"
}
