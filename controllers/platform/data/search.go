package data

import (
	"beegostudy/models"
	"beegostudy/service/bleve"
	"beegostudy/service/progress"
)

type SearchController struct {
	DataController
}

//DataGrid列表数据加载
func (c *SearchController) List() {
	if datagrid, err := models.GetUsersPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
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
	c.put("TaskId", "test")
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
