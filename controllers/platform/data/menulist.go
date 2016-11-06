package data

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
)

type MenuListController struct {
	beego.Controller
}

func (c *MenuListController) Get() {

	c.GetInt("rows")    //每页行数
	c.GetInt("page")    //请求的页码
	c.GetString("sidx") //排序字段
	c.GetString("sord") //升降序

	if menus, err := models.GetMenus(); err != nil {
		beego.Error(err)
	} else {
		jqGrid := new(models.JqGrid)
		jqGrid.Rows = menus
		jqGrid.Page = 1
		jqGrid.Total = 4
		jqGrid.Records = 8
		c.Data["json"] = jqGrid
		c.ServeJSON()
	}
}
