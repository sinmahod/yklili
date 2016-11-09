package data

import (
	"github.com/astaxie/beego"
	"net/http"
)

type DataController struct {
	beego.Controller
	//每页行数
	PageSize int
	//请求的页码
	PageIndex int
	//排序字段
	OrderColumn string
	//升降序
	OrderSord string
	//方法名
	MethodName string
}

/**
*   准备方法，得到页码等信息
**/
func (c *DataController) Prepare() {
	c.PageSize, _ = c.GetInt("rows")
	c.PageIndex, _ = c.GetInt("page")
	c.OrderColumn = c.GetString("sidx")
	c.OrderSord = c.GetString("sord")
	c.MethodName = c.GetString(":method")
}

func (c *DataController) methodNotFind() {
	http.Error(c.Ctx.ResponseWriter, c.MethodName+" Method Not Find", 404)
}
