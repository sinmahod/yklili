package template

import (
	"beegostudy/models"
	"strings"

	"github.com/astaxie/beego"
)

type HTMLController struct {
	beego.Controller
	HTMLPath string
}

/**
*   准备方法，初始化页头、页尾、菜单以及样式和脚本模板
**/
func (c *HTMLController) Prepare() {

	skin := c.Ctx.GetCookie("skin")

	if skin == "" || skin == "null" {
		c.Data["Skin"] = "no-skin"
	} else {
		c.Data["Skin"] = "skin-" + skin
		c.Data["Skin"+skin] = "selected"
	}

	c.HTMLPath = c.GetString(":path")

	//如果使用pjax的请求则只解析局部模板，反之返回解析全部模板
	_, action := c.GetControllerAndAction()

	pjax := c.GetString("_pjax")

	u := c.GetSession("User")

	if strings.EqualFold(action, "Get") && pjax == "" && u != nil {

		defer func() {
			c.Layout = "platform/platform.html"

			c.LayoutSections = make(map[string]string)
			c.LayoutSections["Include"] = "public/include.tpl"
			c.LayoutSections["Script"] = "public/script.tpl"

			user := u.(*models.S_User)
			c.Data["UserName"] = user.GetUserName()
		}()

		if menus, err := models.GetMenusLevel(c.Ctx.Request.RequestURI); err != nil {
			beego.Error(err)
		} else {
			c.Data["Menus"] = menus
		}
	}
}

func (c *HTMLController) Get() {
	c.TplName = c.HTMLPath + ".html"
}
