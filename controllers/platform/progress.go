//进度条任务
package platform

import (
	"beegostudy/service/progress"

	"github.com/astaxie/beego"
)

type ProgController struct {
	beego.Controller
}

func init() {
	beego.Router("/platform/prog", &ProgController{})
}

func (c *ProgController) Post() {
	taskId := c.GetString("taskId")
	c.Data["json"] = progress.GetPerc(taskId)
	c.ServeJSON()
}
