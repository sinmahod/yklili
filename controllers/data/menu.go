package data

import (
	"beegostudy/models"
	"beegostudy/models/orm"
	"beegostudy/util"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

type MenuController struct {
	DataController
}

func (c *MenuController) Get() {
	//得到方法名，利用反射机制获取机构体
	value := reflect.ValueOf(c)
	//判断结构中是否存在方法，存在则执行
	if v := value.MethodByName(c.MethodName); v.IsValid() {
		v.Call(nil)
	} else {
		c.methodNotFind()
	}
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
	idStr := c.GetString("Id")

	if idStr != "" {
		id, _ := strconv.Atoi(idStr)

		menu, err := models.GetMenu(id)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["Menu"] = menu

		menus, err := models.GetTopMenus(menu.GetPid())
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["ParentMenus"] = menus
	}

	c.TplName = "platform/menu/menuDialog.html"
	c.addScript()
}

//保持数据
func (c *MenuController) Save() {
	if len(c.RequestData) > 0 {
		menu := new(models.Menu)
		tran := new(orm.Transaction)
		if util.IsNumber(c.RequestData["Id"]) {
			menu.SetId(c.RequestData["Id"])
			menu.Fill()
		}
		if err := menu.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
			c.fail("操作失败，请确认参数是否传递正确")
		} else {
			if !util.IsNumber(c.RequestData["Id"]) {
				pid := util.Atoi(c.RequestData["Pid"])
				if pid == 0 {
					menu.SetLevel(1)
				} else {
					menu.SetLevel(2)
				}
				menu.SetCurrentTime()
				user := c.GetSession("User").(*models.User)
				menu.SetAddUser(user.GetUserName())
				tran.Add(menu, orm.INSERT)
			} else {
				tran.Add(menu, orm.UPDATE)
			}

			if tran.Commit() != nil {
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

func (c *MenuController) Del() {
	ids := c.GetString("Ids")
	if ids != "" {
		tran := new(orm.Transaction)
		idList := strings.Split(ids, ",")
		for _, id := range idList {
			menu := new(models.Menu)
			menu.SetId(id)
			if !menu.GetIsLeaf() {
				c.fail("操作失败，要删除的菜单存在子级菜单，请先删除子级菜单")
				c.ServeJSON()
				return
			}
			tran.Add(menu, orm.DELETE)
		}
		if tran.Commit() == nil {
			c.success("操作成功")
		} else {
			c.fail("操作失败，传递参数为空")
		}

	} else {
		c.fail("操作失败，传递参数为空")
	}
	c.ServeJSON()
}
