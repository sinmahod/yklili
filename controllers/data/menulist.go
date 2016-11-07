package data

import (
	"beegostudy/models"

	"github.com/astaxie/beego"
)

type MenuListController struct {
	DataController
}

func (c *MenuListController) Get() {
	if datagrid, err := models.GetMenusPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}
