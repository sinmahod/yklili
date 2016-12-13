package data

import (
	"beegostudy/models"
	"beegostudy/service/cron"
	"reflect"
)

type CronController struct {
	DataController
}

func (c *CronController) Get() {
	//得到方法名，利用反射机制获取结构体
	value := reflect.ValueOf(c)
	//判断结构中是否存在方法，存在则执行
	if v := value.MethodByName(c.MethodName); v.IsValid() {
		v.Call(nil)
	} else {
		c.methodNotFind()
	}
}

func (c *CronController) List() {
	tasks := cron.GetTaskList()
	c.Data["json"] = models.GetDataGrid(tasks, c.PageIndex, c.PageSize, int64(len(tasks)))
	c.ServeJSON()
}

func (c *CronController) Run() {
	c.Data["json"] = cron.CronStatus()
	c.ServeJSON()
}

func (c *CronController) Stop() {
	id := c.GetString("TaskId")
	cron.TaskStop(id)
	c.success("操作完成，任务已暂停")
	c.ServeJSON()
}

func (c *CronController) Start() {
	id := c.GetString("TaskId")
	cron.TaskStart(id)
	c.success("操作完成，任务已开始")
	c.ServeJSON()
}

func (c *CronController) Exec() {
	ids := c.GetString("TaskIds")
	cron.TaskExecute(ids)
	c.success("操作完成，任务执行完毕")
	c.ServeJSON()
}
