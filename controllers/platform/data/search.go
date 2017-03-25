package data

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"yklili/models"
	"yklili/models/orm"
	"yklili/service/bleve"
	"yklili/service/progress"
	"yklili/util/numberutil"
)

type SearchController struct {
	DataController
}

//DataGrid列表数据加载
func (c *SearchController) List() {
	if datagrid, err := models.GetWordsPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

//修改/新建初始化
func (c *SearchController) InitPage() {
	idStr := c.GetString("Id")

	if idStr != "" {
		id, _ := strconv.Atoi(idStr)

		words, err := models.GetWordsById(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["Words"] = words
	}

	c.TplName = "platform/search/currentWordsDialog.html"
	c.addScript()
}

func (c *SearchController) VerifySWords() {
	swords := c.GetString("SWords")
	if models.WordsExists(swords) == 0 {
		c.Data["json"] = Validator{true}
	} else {
		c.Data["json"] = Validator{false}
	}
	c.ServeJSON()
}

//保存
func (c *SearchController) Save() {
	if len(c.RequestData) > 0 {
		words := new(models.S_SearchWords)
		tran := new(orm.Transaction)
		if numberutil.IsNumber(c.RequestData["Id"]) {
			words.SetId(c.RequestData["Id"])
			words.Fill()
		}
		if err := words.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
		} else {
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				words.SetId(models.GetMaxId("S_SearchWordsID"))
				words.SetAddUser(sysuser.GetUserName())
				tran.Add(words, orm.INSERT)
			} else {
				words.SetModifyUser(sysuser.GetUserName())
				tran.Add(words, orm.UPDATE)
			}

			if err = tran.Commit(); err != nil {
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

func (c *SearchController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			words := new(models.S_SearchWords)
			words.SetId(id)
			tran.Add(words, orm.DELETE)
		}
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

func (c *SearchController) ImportWords() {
	sysuser := c.GetSession("User").(*models.S_User)
	pt := new(progress.ProgressTask)
	pt.SetTaskId("ImportWords")
	pt.SetFunc(func() {
		err := models.ImportWords(bleve.USER_DICT_PATH, sysuser.GetUserName(), pt)
		if err == nil {
			pt.SetMsg("导入词典完成")
			pt.SetPerc(100)
		} else {
			pt.SetMsg("导入词典失败")
		}
	})
	pt.Start()
	c.put("TaskId", "ImportWords")
	c.success("任务开始")
	c.ServeJSON()
}

//重建索引
func (c *SearchController) RebuildIndex() {

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
	c.put("TaskId", "RebuildIndex")
	c.success("任务开始")
	c.ServeJSON()
}

//索引情况查看
func (c *SearchController) IndexView() {

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
