package data

import (
	"strconv"
	"strings"
	"yklili/models"
	"yklili/models/orm"
	"yklili/util/numberutil"

	"github.com/astaxie/beego"
)

type UserController struct {
	DataController
}

//DataGrid列表数据加载
func (c *UserController) List() {
	if datagrid, err := models.GetUsersPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

//修改/新建初始化
func (c *UserController) InitPage() {
	idStr := c.GetString("Id")

	if idStr != "" {
		id, _ := strconv.Atoi(idStr)

		user, err := models.GetUserById(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["User"] = user
	}

	c.TplName = "platform/user/userDialog.html"
	c.addScript()
}

//修改密码初始化
func (c *UserController) EditPassword() {
	idStr := c.GetString("Id")

	if idStr != "" {
		id, _ := strconv.Atoi(idStr)

		user, err := models.GetUserById(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["User"] = user
	}

	c.TplName = "platform/user/userPasswordDialog.html"
	c.addScript()
}

//保存数据
func (c *UserController) Save() {
	if len(c.RequestData) > 0 {
		user := new(models.S_User)
		tran := new(orm.Transaction)
		if numberutil.IsNumber(c.RequestData["Id"]) {
			user.SetId(c.RequestData["Id"])
			user.Fill()
		}
		if err := user.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
		} else {
			if c.GetString("Password") != "" {
				user.SetPassword(c.GetString("Password"))
			}
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				user.SetId(models.GetMaxId("S_UserID"))
				user.SetAddUser(sysuser.GetUserName())
				tran.Add(user, orm.INSERT)
			} else {
				user.SetModifyUser(sysuser.GetUserName())
				tran.Add(user, orm.UPDATE)
			}

			if err = tran.Commit(); err != nil {
				beego.Error(err)
				c.fail("操作失败，数据修改时出现错误")
			} else {
				c.success("操作成功")
			}
		}
	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}

//修改Miami
func (c *UserController) SavePassword() {
	if len(c.RequestData) > 0 {
		user := new(models.S_User)
		tran := new(orm.Transaction)
		if numberutil.IsNumber(c.RequestData["Id"]) && c.GetString("Password") != "" {
			user.SetId(c.RequestData["Id"])
			user.Fill()
			user.SetPassword(c.GetString("Password"))
			sysuser := c.GetSession("User").(*models.S_User)
			user.SetModifyUser(sysuser.GetUserName())
			tran.Add(user, orm.UPDATE)
			if err := tran.Commit(); err != nil {
				beego.Error(err)
				c.fail("操作失败，数据修改时出现错误")
			} else {
				c.success("操作成功")
			}
		} else {
			beego.Warn("请确认参数是否传递正确")
			c.fail("操作失败，请先选择用户")
		}
	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}

func (c *UserController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			user := new(models.S_User)
			user.SetId(id)
			tran.Add(user, orm.DELETE)
		}
		if err := tran.Commit(); err != nil {
			beego.Error(err)
			c.fail("操作失败，操作数据库时出现错误")
		} else {
			c.success("操作成功")
		}

	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}

//验证用户名是否存在 true不存在,false存在
func (c *UserController) VerifyUserName() {
	username := c.GetString("UserName")
	if models.UserExists(username) == 0 {
		c.Data["json"] = Validator{true}
	} else {
		c.Data["json"] = Validator{false}
	}
	c.ServeJSON()
}

//验证邮箱是否存在 true不存在,false存在
func (c *UserController) VerifyEmail() {
	email := c.GetString("Email")
	if models.EmailExists(email) == 0 {
		c.Data["json"] = Validator{true}
	} else {
		c.Data["json"] = Validator{false}
	}
	c.ServeJSON()
}
