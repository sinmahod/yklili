package routers

import (
	"beegostudy/controllers"
	"fmt"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	fmt.Println("asdasd")
}
