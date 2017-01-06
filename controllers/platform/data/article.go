package data

import (
	"beegostudy/models"
	"beegostudy/models/orm"
	"beegostudy/util/numberutil"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

type ArticleController struct {
	DataController
}

//DataGrid列表数据加载
func (c *ArticleController) List() {
	if datagrid, err := models.GetArticlesPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
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
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				article.SetAddUser(sysuser.GetUserName())
				tran.Add(article, orm.INSERT)
			} else {
				article.SetModifyUser(sysuser.GetUserName())
				tran.Add(article, orm.UPDATE)
			}

			if tran.Commit() != nil {
				beego.Error(err)
				c.fail("操作失败，数据修改时出现错误")
			} else {
				c.success("操作成功")
			}
		}
	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}

func (c *ArticleController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			article := new(models.S_Article)
			article.SetId(id)
			tran.Add(article, orm.DELETE)
		}
		if tran.Commit() == nil {
			c.success("操作成功")
		} else {
			c.fail("操作失败，传递参数为空")
		}

	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}
