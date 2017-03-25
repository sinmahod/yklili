package data

import (
	"strconv"
	"yklili/models"
	"yklili/models/orm"
	"yklili/util/numberutil"

	"github.com/astaxie/beego"
)

type ArticleController struct {
	DataController
}

//DataGrid列表数据加载
func (c *ArticleController) List() {

	var datagrid *models.DataGrid
	var err error

	if c.GetString("type") == "article" {
		datagrid, err = models.GetArticlesPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData)
	} else {
		datagrid, err = models.GetPackages(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData)
	}
	if err != nil {
		beego.Error(err)
	}

	c.Data["json"] = datagrid
	c.ServeJSON()
}

//修改/新建初始化
func (c *ArticleController) InitPage() {
	idStr := c.GetString("Id")

	if idStr != "" {
		id, _ := strconv.Atoi(idStr)

		article, err := models.GetArticle(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["Article"] = article
	}

	c.TplName = "platform/article/articleDetail.html"
	c.addScript()
}

//修改/新建初始化
func (c *ArticleController) InitArticle() {
	idStr := c.GetString("Id")

	if idStr != "" {
		id, _ := strconv.Atoi(idStr)

		article, err := models.GetArticle(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.put("Article", article)
		c.success("操作成功")
	} else {
		c.fail("请先选择文章")
	}
	c.ServeJSON()
}

//保持数据
func (c *ArticleController) Save() {
	if len(c.RequestData) > 0 {
		article := new(models.S_Article)
		tran := new(orm.Transaction)
		if numberutil.IsNumber(c.RequestData["Id"]) {
			article.SetId(c.RequestData["Id"])
			article.Fill()
		}
		if err := article.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
		} else {
			if numberutil.IsNumber(c.RequestData["Status"]) {
				article.SetStatus(models.PUBLISH)
			} else {
				article.SetStatus(models.DRAFT)
			}
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				article.SetId(models.GetMaxId("S_ArticleID"))
				article.SetAddUser(sysuser.GetUserName())
				tran.Add(article, orm.INSERT)
			} else {
				article.SetModifyUser(sysuser.GetUserName())
				tran.Add(article, orm.UPDATE)
			}

			if err = tran.Commit(); err != nil {
				beego.Error(err)
				c.fail("操作失败，数据修改时出现错误")
			} else {
				c.put("Id", article.GetId())
				c.success("操作成功")
			}
		}
	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}

func (c *ArticleController) Del() {
	id := c.GetString("Id")
	if id != "" {
		tran := new(orm.Transaction)
		article := new(models.S_Article)
		article.SetId(id)
		article.SetStatus(models.DELETE)
		tran.Add(article, orm.UPDATE)
		if err := tran.Commit(); err != nil {
			beego.Error(err)
			c.fail("操作失败，操作数据库时出现错误")
		} else {
			c.success("操作成功")
		}
	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}
