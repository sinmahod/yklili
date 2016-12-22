package data

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
)

type ImageController struct {
	DataController
}

//DataGrid列表数据加载
func (c *ImageController) List() {
	c.RequestData["filetype"] = "image/*"
	if datagrid, err := models.GetAttchmentsPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}
