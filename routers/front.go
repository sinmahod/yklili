package routers

import (
	"fmt"
	"github.com/astaxie/beego"
	"yklili/controllers/front"
)

func init() {
	//数据控制器
	models := map[string]beego.ControllerInterface{
		"article": &front.ArticleController{},
		"catalog": &front.CatalogController{},
		"search":  &front.SearchController{},
	}

	for name, controller := range models {
		beego.Router(fmt.Sprintf("/%s", name), controller, "*:Page")
		beego.Router(fmt.Sprintf("/%s/:method", name), controller, "*:Get")
	}
	beego.Router("/", &front.IndexController{})
}
