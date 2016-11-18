package dml

import "github.com/astaxie/beego"

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

func (c *DMLController) Get() {
	c.TplName = c.DMLPath + ".dml"
}
