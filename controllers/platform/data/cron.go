package data

import (
	"beegostudy/models"
	"beegostudy/service/cron"
)

type CronController struct {
	DataController
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
