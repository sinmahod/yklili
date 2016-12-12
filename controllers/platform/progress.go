package platform

import "github.com/astaxie/beego"

type ProgController struct {
	beego.Controller
}

func (c *ProgController) Page() {
	c.ServeJSON()
}

func (c *ProgController) Post() {
	//id := c.GetString("Id")
	// if t, ok := threadtask[id]; ok {
	// 	c.Data["json"] = t.GetPerc()
	// } else {
	// 	c.Data["json"] = ok
	// }
	c.Data["json"] = 123
	c.ServeJSON()
}
