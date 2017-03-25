package data

import (
	"time"
	"yklili/service/progress"
)

type TestController struct {
	DataController
}

func (c *TestController) Exec() {
	pt := new(progress.ProgressTask)
	pt.SetTaskId("test")
	pt.SetFunc(func() {
		for i := 0; i < 100; i++ {
			time.Sleep(100000000)
			pt.SetPerc(i)
			pt.SetMsg("任务已执行到了%d%s", i, "%")
		}
	})
	pt.Start()
	c.put("TaskId", "test")
	c.success("任务开始")
	c.ServeJSON()
}
