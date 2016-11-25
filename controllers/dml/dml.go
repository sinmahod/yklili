package dml

import (
	"bytes"
	"github.com/astaxie/beego"
)

type DMLController struct {
	beego.Controller
	DMLPath string
}

/**
*   准备方法，得到页面URL，专门处理dml文件（Dialog文件）
**/
func (c *DMLController) Prepare() {
	c.DMLPath = c.GetString(":path")
}

func (c *DMLController) Finish() {
	c.TplName = c.DMLPath + ".dml"
}

func (c *DMLController) Render() error {
	if !c.EnableRender {
		return nil
	}
	rb, err := c.RenderBytes()
	if err != nil {
		return err
	}
	html := AddScript(string(rb))
	rb = []byte(html)
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	return c.Ctx.Output.Body(rb)
}

const (
	Script string = "<script>$('[data-rel=tooltip]').tooltip({container:'body'});</script>"
)

func AddScript(content string) string {
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲
	buffer.WriteString(content)
	buffer.WriteString("\n")
	buffer.WriteString(Script)
	return buffer.String()
}
