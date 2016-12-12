package data

import (
	"beegostudy/service/progress"
	"fmt"
	"reflect"
	"time"
)

type TestController struct {
	DataController
}

func (c *TestController) Get() {
	//得到方法名，利用反射机制获取机构体
	value := reflect.ValueOf(c)
	//判断结构中是否存在方法，存在则执行
	if v := value.MethodByName(c.MethodName); v.IsValid() {
		v.Call(nil)
	} else {
		c.methodNotFind()
	}
}

func (c *TestController) Exec() {
	Test()
	c.success("操作完成，任务执行完毕")
	c.ServeJSON()
}

func Test() {
	pt := new(progress.ProgressTask)
	pt.SetTaskId("test")
	pt.SetFunc(func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Second * 1)
			fmt.Println(i)
			pt.SetPerc(i)
		}
	})
	fmt.Println(pt.Start())
}
