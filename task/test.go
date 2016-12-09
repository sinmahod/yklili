package task

import (
	"beegostudy/service/cron"
	"fmt"
)

func init() {
	cron.RegisterTask(&TestTask{})
}

type TestTask struct{}

func (t *TestTask) GetId() string {
	return "Test"
}

func (t *TestTask) GetSpec() string {
	return "*/3  *  *  *  *  *"
}

func (t *TestTask) GetDesc() string {
	return "测试任务（每3秒执行一次）"
}

func (t *TestTask) Execute() {
	fmt.Println("------------测试任务（每3秒执行一次）-------------")
}
