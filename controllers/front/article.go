package front

import (
	"github.com/astaxie/beego"
	"github.com/sinmahod/yklili/models"
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

func (c *ArticleController) PackageList() {
	if datagrid, err := models.GetFrontPackages(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

func (c *ArticleController) Page() {
	if id, err := c.GetInt("id"); err != nil {
		beego.Error(err)
	} else {
		article, err := models.GetArticleByStatus(id, 1)
		if err != nil {
			beego.Error(err)
		} else {
			c.Data["Article"] = article
		}
	}
	c.TplName = "front/article.html"
}
