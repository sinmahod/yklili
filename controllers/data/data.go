package data

import (
	"bytes"
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
	//存放字段和值的数据
	RequestData map[string]interface{}
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

	if err := c.Ctx.Request.ParseForm(); err != nil {
		beego.Info(err)
	}
	c.RequestData = make(map[string]interface{})
	for k, v := range c.Ctx.Request.Form {
		if len(v) > 0 {
			c.RequestData[k] = v[0]
		}
	}
}

func (c *DataController) methodNotFind() {
	http.Error(c.Ctx.ResponseWriter, c.MethodName+" 方法未找到", 404)
}

func (c *DataController) paramIsNull() {
	http.Error(c.Ctx.ResponseWriter, "参数为空", 510)
}

const (
	Script string = "<script>$('[data-rel=tooltip]').tooltip({container:'body'});</script>"
)

//在模板结尾追加js
func (c *DataController) AddScript() error {
	rb, err := c.RenderBytes()
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	buffer.WriteString(string(rb))
	buffer.WriteString("\n")
	buffer.WriteString(Script)

	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	c.Ctx.Output.Body(buffer.Bytes())
	return nil
}
