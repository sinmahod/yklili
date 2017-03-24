package cache

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

func InitRedis() error {
	redisIP := beego.AppConfig.String("Redis.IP")
	redisPort := beego.AppConfig.String("Redis.Port")
	redisPassWord := beego.AppConfig.String("Redis.Password")
	redisKey := beego.AppConfig.String("Redis.Key")
	redisDBNum := beego.AppConfig.String("Redis.DBNum")

	redisInfo := make(map[string]string)
	redisInfo["conn"] = redisIP + ":" + redisPort
	redisInfo["password"] = redisPassWord
	redisInfo["key"] = redisKey
	redisInfo["dbNum"] = redisDBNum

	redisJson, _ := json.Marshal(redisInfo)

	var err error
	c, err = cache.NewCache("redis", string(redisJson))
	if err != nil {
		return err
	}
	beego.Debug("缓存连接信息 : ", string(redisJson))
	return nil
}
