package data

import (
	"github.com/astaxie/beego"
	"github.com/sinmahod/yklili/models"
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

func (c *ImageController) InitPage() {
	c.TplName = "platform/image/imageDialog.html"
	c.addScript()
}
