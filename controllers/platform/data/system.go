package data

import (
	"beegostudy/models"
	"github.com/astaxie/beego"
)

type SystemController struct {
	DataController
}

func (c *SystemController) List() {
	m := make(map[string]interface{})
	m["MaxMemory"] = beego.BConfig.MaxMemory           //最大内存
	m["ServerName"] = beego.BConfig.ServerName         //Beego版本
	m["AppPath"] = beego.AppPath                       //应用所在路径
	m["AppName"] = beego.BConfig.AppName               //应用名称
	m["HTTPPort"] = beego.BConfig.Listen.HTTPPort      //端口
	m["RunMode"] = beego.BConfig.RunMode               //运行模式
	m["StaticDir"] = beego.BConfig.WebConfig.StaticDir //静态目录
	c.Data["json"] = models.GetDataGrid(m, 0, 0, 7)
	c.ServeJSON()
}
