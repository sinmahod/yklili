package data

import (
	"github.com/sinmahod/yklili/models"
	"github.com/sinmahod/yklili/models/orm"
	"github.com/sinmahod/yklili/util/numberutil"
	"strings"

	"github.com/astaxie/beego"
)

type PackageController struct {
	DataController
}

//修改/新建初始化
func (c *PackageController) InitPage() {
	if numberutil.IsNumber(c.RequestData["Id"]) {
		id := numberutil.Atoi(c.RequestData["Id"])

		p, err := models.GetPackage(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["Package"] = p
	}
	c.TplName = "platform/package/packageDialog.html"
	c.addScript()
}

//保存数据
func (c *PackageController) Save() {
	if len(c.RequestData) > 0 {
		p := new(models.S_Package)
		tran := new(orm.Transaction)

		if numberutil.IsNumber(c.RequestData["Id"]) {
			p.SetId(c.RequestData["Id"])
			p.Fill()
		}

		if err := p.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
		} else {
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				p.SetId(models.GetMaxId("S_PackageID"))
				p.SetAddUser(sysuser.GetUserName())
				tran.Add(p, orm.INSERT)
			} else {
				p.SetModifyUser(sysuser.GetUserName())
				tran.Add(p, orm.UPDATE)
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

func (c *PackageController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			p := new(models.S_Package)
			p.SetId(id)
			if !p.IsNull() {
				c.fail("操作失败，要删除的文件夹内存在多篇文章，请先删除文章")
				c.ServeJSON()
				return
			}
			tran.Add(p, orm.DELETE)
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
