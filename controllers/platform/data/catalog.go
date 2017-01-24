package data

import (
	"beegostudy/models"
	"beegostudy/models/orm"
	"beegostudy/util/numberutil"
	"github.com/astaxie/beego"
	"strings"
)

type CatalogController struct {
	DataController
}

//DataGrid列表数据加载
func (c *CatalogController) List() {
	if datagrid, err := models.GetCatalogsPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

//修改/新建初始化
func (c *CatalogController) InitPage() {
	var cs []models.CatalogSelectInit

	if numberutil.IsNumber(c.RequestData["Id"]) {
		id := numberutil.Atoi(c.RequestData["Id"])

		catalog, err := models.GetCatalog(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["Catalog"] = catalog

		cs, _ = models.GetTopCatalogs(catalog.GetPid(), id)
	} else {
		cs, _ = models.GetTopCatalogs(0, 0)
	}
	c.Data["ParentCatalogs"] = cs
	c.TplName = "platform/catalog/catalogDialog.html"
	c.addScript()
}

//保存数据
func (c *CatalogController) Save() {
	if len(c.RequestData) > 0 {
		catalog := new(models.S_Catalog)
		tran := new(orm.Transaction)
		pid := numberutil.Atoi(c.RequestData["Pid"])

		isNewParent := false

		if numberutil.IsNumber(c.RequestData["Id"]) {
			catalog.SetId(c.RequestData["Id"])
			catalog.Fill()
			isNewParent = pid != catalog.GetPid()
		}

		if err := catalog.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
			goto END
		} else {
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				if pid == 0 {
					catalog.SetLevel(1)
					catalog.SetInnerCode(models.GetMaxNo("catalog", "", 4))
				} else {
					catalog.SetLevel(2)
					catalog.SetInnerCode(models.GetMaxNo("catalog", models.GetCatalogInnerCode(pid), 4))
				}
				catalog.SetId(models.GetMaxId("S_CatalogID"))
				catalog.SetAddUser(sysuser.GetUserName())
				tran.Add(catalog, orm.INSERT)
			} else {
				if isNewParent {
					if pid != 0 {
						//如果不是叶子节点则不允许改变父级ID
						if !catalog.GetIsLeaf() {
							c.fail("操作失败，当前菜单存在子级菜单，请先清空子级菜单")
							goto END
						}
						catalog.SetLevel(2)
						pcode := models.GetCatalogInnerCode(pid)
						catalog.SetInnerCode(models.GetMaxNo("catalog", pcode, 4))
					} else {
						catalog.SetLevel(1)
						catalog.SetInnerCode(models.GetMaxNo("catalog", "", 4))
					}
				}
				catalog.SetModifyUser(sysuser.GetUserName())
				tran.Add(catalog, orm.UPDATE)
			}

			//if err = tran.Commit(); err != nil {
			beego.Error(err)
			c.fail("操作失败，数据修改时出现错误")
			//} else {
			c.success("操作成功")
			//}
		}
	} else {
		c.fail("操作失败，传递参数为空")
	}
END:
	c.ServeJSON()
}

func (c *CatalogController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			catalog := new(models.S_Catalog)
			catalog.SetId(id)
			if !catalog.GetIsLeaf() {
				c.fail("操作失败，要删除的栏目存在子级栏目，请先删除子级栏目")
				c.ServeJSON()
				return
			}
			tran.Add(catalog, orm.DELETE)
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
