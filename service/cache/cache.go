package cache

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
)

var c cache.Cache

func init() {
	cacheState := beego.AppConfig.String("Cache")
	if cacheState == "true" {
		beego.Info("正在初始化缓存连接...")
		defer func() {
			if r := recover(); r != nil {
				beego.Error("initial cache error caught: %v\n", r)
				c = nil
			}
		}()

		cacheType := beego.AppConfig.String("Cache.Type")

		if cacheType == "redis" {
			if err := InitRedis(); err != nil {
				beego.Error("缓存创建失败：", err)
				return
			}

		}

		beego.Info("缓存连接成功")
	}
}
