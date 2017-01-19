package conf

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
	"strconv"
)

var configItem map[string]string

func config(in string) string {
	return GetValue(in)
}

func init() {
	RegisterDB()
	beego.AddFuncMap("SiteConfig", config)
	configItem = make(map[string]string)
	Reload()
}

// 刷新配置
func Reload() {
	Clear()
	for _, c := range models.GetConfigs() {
		configItem[c.GetK()] = c.GetV()
	}
	if siteid := beego.AppConfig.String("siteid"); siteid != "" {
		id, err := strconv.Atoi(siteid)
		site, err := models.GetSite(id)
		if err == nil {
			configItem["Banner"] = site.GetBanner()
			configItem["Name"] = site.GetName()
			configItem["Subtitle"] = site.GetSubtitle()
			configItem["Host"] = site.GetHost()
			configItem["BackgroundColor"] = site.GetBackgroundColor()
			configItem["Copyright"] = site.GetCopyright()
			configItem["Desc"] = site.GetDesc()
			configItem["ArticleColor"] = site.GetArticleColor()
			return
		}
	}
	configItem["Banner"] = "/static/front/images/banner.jpg"
	configItem["Name"] = "请前往后台设置站点"
	configItem["Subtitle"] = "请前往后台设置站点"
	configItem["Host"] = "/"
	configItem["BackgroundColor"] = "#2c2c2c"
	configItem["Desc"] = "请前往后台设置站点"
	configItem["Copyright"] = "&copy; 2017 All rights  Reserved"
	configItem["ArticleColor"] = "#fff"
}

// 清空配置
func Clear() {
	for k, _ := range configItem {
		delete(configItem, k)
	}
}

// 保存配置
func Save() {

}

// 移除配置
func Remove(k string) {
	delete(configItem, k)
}

// 添加配置
func Add(k, v string) {
	configItem[k] = v
}

const (
	//文件上传路径
	UploadPath = "UploadPath"
)

// 获取配置
func GetValue(k string) string {
	if v, ok := configItem[k]; ok {
		return v
	}
	return ""
}
