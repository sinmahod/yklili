package data

import (
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

func (c *TestController) Stat() {
	id := c.GetString("Id")
	if t, ok := threadtask[id]; ok {
		c.Data["json"] = t.GetPerc()
	} else {
		c.Data["json"] = ok
	}

	c.ServeJSON()
}

var threadtask = make(map[string]*PercentageTask)

type PercentageTask struct {
	id   string
	perc int
	fc   func()
}

func (t *PercentageTask) SetTaskId(taskId string) error {
	t.id = taskId
	return nil
}

func (t *PercentageTask) Start() error {
	if tt, ok := threadtask[t.id]; ok {
		return fmt.Errorf("任务[%s]已经在执行中，请等待完成，当前进度为(%d%s)", t.id, tt.perc, "%")
	}
	go t.fc()
	threadtask[t.id] = t
	return nil

}

func (t *PercentageTask) SetFunc(fc func()) {
	t.fc = fc
}

// 得到任务当前进度
func (t *PercentageTask) GetPerc() int {
	return t.perc
}

// 得到任务当前进度
func (t *PercentageTask) SetPerc(perc int) {
	t.perc = perc
}

func Test() {
	pt := new(PercentageTask)
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
