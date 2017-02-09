package data

import (
	"beegostudy/models"
	"beegostudy/models/orm"
	"beegostudy/service/progress"
	"beegostudy/util/numberutil"
	"strconv"

	"beegostudy/service/bleve"
	"github.com/astaxie/beego"
)

type ArticleController struct {
	DataController
}

//重建索引
func (c *ArticleController) RebuildIndex() {

	pt := new(progress.ProgressTask)
	pt.SetTaskId("RebuildIndex")
	pt.SetFunc(func() {
		err := models.RebuildIndex(pt)
		if err == nil {
			pt.SetMsg("索引重建完成")
			pt.SetPerc(100)
		} else {
			pt.SetMsg("索引创建失败")
		}
	})
	pt.Start()
	c.put("TaskId", "test")
	c.success("任务开始")
	c.ServeJSON()
}

//索引情况查看
func (c *ArticleController) IndexView() {

	q := c.GetString("q")

	p, _ := c.GetInt("p")

	if q != "" {
		if p == 0 {
			p++
		}

		data, err := bleve.And(q).SearchToData(10, p)
		if err != nil {
			c.put("Result", err)
		} else {
			c.put("Result", data)
		}
	}

	cnt, _ := bleve.GetDocCount()

	c.put("Count", cnt)
	c.success("完成")
	c.ServeJSON()
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
