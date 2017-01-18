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
			mp["subtitle"] = site.GetSiteSubtitle()
			mp["host"] = site.GetHost()
			mp["desc"] = site.GetDesc()
			mp["backgroundcolor"] = site.GetBackgroundColor()
			mp["copyright"] = site.GetCopyright()
			c.Data["Index"] = mp
			return
		}
	}

	mp["banner"] = "/static/front/images/banner.jpg"
	mp["title"] = "请前往后台设置站点"
	mp["subtitle"] = "请前往后台设置站点"
	mp["host"] = "/"
	mp["backgroundcolor"] = "#2c2c2c"
	mp["desc"] = "请前往后台设置站点"
	mp["copyright"] = "&copy; 2017 All rights  Reserved"
	c.Data["Index"] = mp
}
