package data

import (
	"github.com/sinmahod/yklili/conf"
	"github.com/sinmahod/yklili/models"
	"github.com/sinmahod/yklili/models/orm"
	"github.com/sinmahod/yklili/util/numberutil"
	"strings"

	"github.com/astaxie/beego"
)

type SiteController struct {
	DataController
}

//DataGrid列表数据加载
func (c *SiteController) List() {
	if datagrid, err := models.GetSitePage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

//修改/新建初始化
func (c *SiteController) InitPage() {
	site, err := models.GetSite()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = site
	}

	c.ServeJSON()
}

//保存数据
func (c *SiteController) Save() {
	if len(c.RequestData) > 0 {
		site := new(models.S_Site)
		tran := new(orm.Transaction)
		if numberutil.IsNumber(c.RequestData["Id"]) {
			site.SetId(c.RequestData["Id"])
			site.Fill()
		}
		if err := site.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
		} else {
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				site.SetId(models.GetMaxId("S_SiteID"))
				site.SetAddUser(sysuser.GetUserName())
				tran.Add(site, orm.INSERT)
			} else {
				site.SetModifyUser(sysuser.GetUserName())
				tran.Add(site, orm.UPDATE)
			}

			if err = tran.Commit(); err != nil {
				beego.Error(err)
				c.fail("操作失败，数据修改时出现错误")
			} else {
				conf.Reload()
				c.success("操作成功")
			}
		}
	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}

func (c *SiteController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			site := new(models.S_Site)
			site.SetId(id)
			tran.Add(site, orm.DELETE)
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
