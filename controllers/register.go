package controllers

import (
	"fmt"

	"beegostudy/models"
	"beegostudy/util"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	c.TplName = "login.html"
}

type result struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"errormsg"`
	Link     string `json:"link"`
}

func (c *RegisterController) Post() {
	u := c.GetString("registerusername")
	p := c.GetString("registerpassword")
	e := c.GetString("email")
	phone := c.GetString("phone")

	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	user := new(models.User)
	user.UserName = u
	err := o.Read(user, "UserName")
	fmt.Println("err:", err)
	if err == nil {
		jsondata := &result{201, "您的用户名已存在，请更换用户名！", "./test"}
		//js, _ := json.Marshal(jsondata)

		c.Data["json"] = jsondata
		c.ServeJSON()
		return
	}
	user.Email = e
	errEmail := o.Read(user, "Email")
	fmt.Println("errEmail:", errEmail)
	if errEmail == nil {
		jsondata := &result{201, "您填写的邮箱已存在，请更换邮箱！", "./test"}
		//js, _ := json.Marshal(jsondata)

		c.Data["json"] = jsondata
		c.ServeJSON()
		return
	}
	user.Password = util.GeneratePWD(p)
	user.Phone = phone
	user.AddTime = time.Now()
	user.AddUser = u

	fmt.Println(o.Insert(user))

}
