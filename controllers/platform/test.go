package platform

type TestController struct {
	PlatformController
}

func (c *TestController) Get() {
	c.TplName = "test2.html"
}
