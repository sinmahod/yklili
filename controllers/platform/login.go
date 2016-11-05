package platform

import (
	"beegostudy/models"
	"beegostudy/util"
	"fmt"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	u := c.GetString("username")
	p := c.GetString("password")

	var jsondata result

	if u == "" {
		jsondata = result{0, "请输入用户名！", ""}
		c.Data["json"] = jsondata
		c.ServeJSON()
		return
	}

	if p == "" {
		jsondata = result{0, "请输入密码！", ""}
		c.Data["json"] = jsondata
		c.ServeJSON()
		return
	}

	user, err := models.GetUser(u)
	if err != nil {
		//用户不存在或者读取数据的错误
		jsondata = result{0, fmt.Sprintf("%s", err), ""}
	} else {
		//用户存在则校验用户密码是否正确
		if util.VerifyPWD(p, user.GetPassword()) {
			jsondata = result{1, "", "./platform/users"}
			c.SetSession("User", user)
		} else {
			jsondata = result{0, "密码错误请重试！", ""}
		}
	}
	c.Data["json"] = jsondata
	c.ServeJSON()
}
