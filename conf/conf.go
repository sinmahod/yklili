package conf

import (
	"runtime"

	"github.com/astaxie/beego"
	"github.com/sinmahod/yklili/models"
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

	// 设置默认的资源上传目录
	if _, ok := configItem[UploadPath]; !ok {
		configItem[UploadPath] = beego.AppPath + "/upload/"
	}

	site, err := models.GetSite()
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
	} else {
		configItem["Banner"] = "/static/front/images/banner.jpg"
		configItem["Name"] = "请前往后台设置站点"
		configItem["Subtitle"] = "请前往后台设置站点"
		configItem["Host"] = "/"
		configItem["BackgroundColor"] = "#2c2c2c"
		configItem["Desc"] = "请前往后台设置站点"
		configItem["Copyright"] = "&copy; 2017 All rights  Reserved"
		configItem["ArticleColor"] = "#fff"
	}

}

func GetCurrentFilePath() string {
	_, filePath, _, _ := runtime.Caller(1)
	return filePath
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

func GetDefaultImage() string {
	return beego.AppPath + "/static/images/notfound.gif"
}

// 获取配置
func GetValue(k string) string {
	if v, ok := configItem[k]; ok {
		return v
	}
	return ""
}
