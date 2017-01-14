package front

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
)

type ArticleController struct {
	FrontController
}

func (c *ArticleController) List() {
	if datagrid, err := models.GetArticlesPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

func (c *ArticleController) Page() {
	c.TplName = "front/detail.html"
}
