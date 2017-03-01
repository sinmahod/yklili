package front

import (
	"beegostudy/models"
	"beegostudy/service/bleve"
)

type SearchController struct {
	FrontController
}

func (c *SearchController) Page() {
	c.TplName = "front/search.html"
}

func (c *SearchController) Search() {
	q := c.GetString("q")

	p, _ := c.GetInt("p")

	var as []models.S_Article

	var size = 10

	if q != "" {
		if p == 0 {
			p++
		}

		cnt, err := bleve.And(q).Search(size, p, &as)
		if err != nil {
			panic(err)
		} else {
			pagetotal := cnt / uint64(size)

			if cnt%uint64(size) > 0 {
				pagetotal++
			}

			datagrid := models.GetDataGrid(as, p, int(pagetotal), int64(cnt))
			c.Data["json"] = datagrid
		}
	}

	c.ServeJSON()
}
