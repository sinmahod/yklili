package conf

import (
	"beegostudy/models"
)

var configItem map[string]string

func init() {
	configItem = make(map[string]string)
	//Reload()
}

// 刷新配置
func Reload() {
	Clear()
	for _, c := range models.GetConfigs() {
		configItem[c.GetK()] = c.GetV()
	}
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
