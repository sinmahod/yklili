package front

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
)

type CatalogController struct {
	FrontController
}

func (c *CatalogController) List() {
	if datagrid, err := models.GetCatalogsPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

func (c *CatalogController) Page() {}
