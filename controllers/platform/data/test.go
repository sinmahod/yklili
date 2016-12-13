package data

import (
	"beegostudy/service/progress"
	"reflect"
	"time"
)

type TestController struct {
	DataController
}

func (c *TestController) Get() {
	//得到方法名，利用反射机制获取结构体
	value := reflect.ValueOf(c)
	//判断结构中是否存在方法，存在则执行
	if v := value.MethodByName(c.MethodName); v.IsValid() {
		v.Call(nil)
	} else {
		c.methodNotFind()
	}
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
