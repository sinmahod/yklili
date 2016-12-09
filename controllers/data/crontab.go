package data

import (
	"beegostudy/service/cron"
	"fmt"
	"reflect"
)

type CronController struct {
	DataController
}

func (c *CronController) Get() {
	//得到方法名，利用反射机制获取机构体
	value := reflect.ValueOf(c)
	//判断结构中是否存在方法，存在则执行
	if v := value.MethodByName(c.MethodName); v.IsValid() {
		v.Call(nil)
	} else {
		c.methodNotFind()
	}
}

func (c *CronController) List() {
	c.Data["json"] = cron.GetTaskList()
	c.ServeJSON()
}

func (c *CronController) Run() {

	c.Data["json"] = cron.CronStatus()
	c.ServeJSON()
}

func (c *CronController) Stop() {
	id := c.GetString("TaskId")
	fmt.Println("Stop:", id)
	cron.TaskStop(id)
	c.Data["json"] = cron.TaskStatus(id)
	c.ServeJSON()
}

func (c *CronController) Start() {
	id := c.GetString("TaskId")
	cron.TaskStart(id)
	c.Data["json"] = cron.TaskStatus(id)
	c.ServeJSON()
}
