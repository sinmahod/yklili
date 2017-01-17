package front

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
	"strconv"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.TplName = "front/index.html"

	mp := make(map[string]string)

	if siteid := beego.AppConfig.String("siteid"); siteid != "" {
		id, err := strconv.Atoi(siteid)
		site, err := models.GetSite(id)
		if err == nil {
			mp["banner"] = site.GetBanner()
			mp["title"] = site.GetName()
			mp["host"] = site.GetHost()
			c.Data["Index"] = mp
			return
		}
	}

	mp["banner"] = "/static/front/images/banner.jpg"
	mp["title"] = beego.AppConfig.String("appname")
	c.Data["Index"] = mp
}
