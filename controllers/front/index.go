package front

import ()

type IndexController struct {
	FrontController
}

func (c *IndexController) Get() {
	c.TplName = "front/index.html"

}
