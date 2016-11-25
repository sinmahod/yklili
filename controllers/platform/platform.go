package platform

import (
	"beegostudy/models"

	"github.com/astaxie/beego"
)

type PlatformController struct {
	beego.Controller
}

/**
*   准备方法，初始化页头、页尾、菜单以及样式和脚本模板
**/
func (c *PlatformController) Prepare() {
	//如果使用pjax的请求则只解析局部模板，反之返回解析全部模板
	pjax := c.GetString("_pjax")

	if pjax == "" {

		defer func() {
			c.Layout = "platform/platform.html"

			c.LayoutSections = make(map[string]string)
			c.LayoutSections["Include"] = "public/include.tpl"
			c.LayoutSections["Script"] = "public/script.tpl"

			user := c.GetSession("User").(*models.User)
			c.Data["UserName"] = user.GetUserName()
		}()

		if menus, err := models.GetMenusLevel(c.Ctx.Request.RequestURI); err != nil {
			beego.Error(err)
		} else {
			c.Data["Menus"] = menus
		}
	}
}
