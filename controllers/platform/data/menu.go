package data

import (
	"beegostudy/models"
	"beegostudy/models/orm"
	"beegostudy/util/numberutil"
	"strings"

	"github.com/astaxie/beego"
)

type MenuController struct {
	DataController
}

//DataGrid列表数据加载
func (c *MenuController) List() {
	if datagrid, err := models.GetMenusPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}

//修改/新建初始化
func (c *MenuController) InitPage() {
	var menus []models.MenuSelectInit

	if numberutil.IsNumber(c.RequestData["Id"]) {
		id := numberutil.Atoi(c.RequestData["Id"])

		menu, err := models.GetMenu(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["Menu"] = menu

		menus, _ = models.GetTopMenus(menu.GetPid(), id)
	} else {
		menus, _ = models.GetTopMenus(0, 0)
	}
	c.Data["ParentMenus"] = menus
	c.TplName = "platform/menu/menuDialog.html"
	c.addScript()
}

//保存数据
func (c *MenuController) Save() {
	if len(c.RequestData) > 0 {
		menu := new(models.S_Menu)
		tran := new(orm.Transaction)
		pid := numberutil.Atoi(c.RequestData["Pid"])

		isNewParent := false

		if numberutil.IsNumber(c.RequestData["Id"]) {
			menu.SetId(c.RequestData["Id"])
			menu.Fill()
			isNewParent = pid != menu.GetPid()
		}

		if err := menu.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
			goto END
		} else {
			sysuser := c.GetSession("User").(*models.S_User)
			if !numberutil.IsNumber(c.RequestData["Id"]) {
				if pid == 0 {
					menu.SetLevel(1)
					menu.SetInnerCode(models.GetMaxNo("menu", "", 4))
				} else {
					menu.SetLevel(2)
					menu.SetInnerCode(models.GetMaxNo("menu", models.GetInnerCode(pid), 4))
				}
				menu.SetId(models.GetMaxId("S_MenuID"))
				menu.SetAddUser(sysuser.GetUserName())
				tran.Add(menu, orm.INSERT)
			} else {
				if isNewParent {
					if pid != 0 {
						//如果不是叶子节点则不允许改变父级ID
						if !menu.GetIsLeaf() {
							c.fail("操作失败，当前菜单存在子级菜单，请先清空子级菜单")
							goto END
						}
						menu.SetLevel(2)
						pcode := models.GetInnerCode(pid)
						menu.SetInnerCode(models.GetMaxNo("menu", pcode, 4))
					} else {
						menu.SetLevel(1)
						menu.SetInnerCode(models.GetMaxNo("menu", "", 4))
					}
				}
				menu.SetModifyUser(sysuser.GetUserName())
				tran.Add(menu, orm.UPDATE)
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
END:
	c.ServeJSON()
}

func (c *MenuController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			menu := new(models.S_Menu)
			menu.SetId(id)
			if !menu.GetIsLeaf() {
				c.fail("操作失败，要删除的菜单存在子级菜单，请先删除子级菜单")
				c.ServeJSON()
				return
			}
			tran.Add(menu, orm.DELETE)
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
