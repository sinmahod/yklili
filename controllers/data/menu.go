package data

import (
	"beegostudy/models"
	"reflect"
	"strconv"

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

		menus, err := models.GetTopMenus(menu.Pid)
		if err != nil {
			beego.Error(err)
			return
		}
		c.Data["ParentMenus"] = menus
	}

	c.TplName = "platform/menu/menuDialog.html"
	c.AddScript()
}

//保持数据
func (c *MenuController) Save() {
	if len(c.RequestData) > 0 {
		for k, v := range c.RequestData {
			beego.Info("=======", k, v)
		}
		menu := new(models.Menu)
		if err := menu.SetValue(c.RequestData); err != nil {
			beego.Warn("请确认参数是否传递正确", err)
		} else {
			beego.Info(menu)
		}
	} else {
		c.paramIsNull()
	}
}
