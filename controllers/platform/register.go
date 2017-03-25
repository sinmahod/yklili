package platform

import (
	"github.com/sinmahod/yklili/models"
	"github.com/sinmahod/yklili/util/pwdutil"

	"github.com/astaxie/beego"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	c.TplName = "login.html"
}

func (c *RegisterController) Post() {
	u := c.GetString("username")
	p := c.GetString("password")
	e := c.GetString("email")
	phone := c.GetString("phone")

	var jsondata result

	if _, err := models.GetUser(u); err == nil {
		jsondata = result{0, "您的用户名已存在，请更换用户名！", ""}
		c.Data["json"] = jsondata
		c.ServeJSON()
		return
	}
	if _, err := models.GetUserByEmail(e); err == nil {
		jsondata = result{0, "您的邮箱已存在，请更换邮箱！", ""}
		c.Data["json"] = jsondata
		c.ServeJSON()
		return
	}

	if _, err := models.InsertUser(u, pwdutil.GeneratePWD(p), e, phone); err != nil {
		jsondata = result{0, "数据库操作失败！", ""}
	} else {
		jsondata = result{1, "注册成功", ""}
	}
	c.Data["json"] = jsondata
	c.ServeJSON()
	return
}
