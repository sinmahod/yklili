package controllers

import (
	"fmt"

	//"beegostudy/models"

	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

type data struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"errormsg"`
	Link     string `json:"link"`
}

func (c *LoginController) Post() {
	u := c.GetString("username")
	p := c.GetString("password")
	fmt.Println(u, p)
	if u != "" {
		jsondata := &data{201, "您的用户名或密码输入错误，请重试！", "./test"}
		//js, _ := json.Marshal(jsondata)

		c.Data["json"] = jsondata
		c.ServeJSON()
	} else {

	}
}
