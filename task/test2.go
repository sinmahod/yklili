package task

import (
	"beegostudy/service/cron"
	"fmt"
)

func init() {
	cron.RegisterTask(&TestTaskA{})
}

type TestTaskA struct{}

func (t *TestTaskA) GetId() string {
	return "TestA"
}

func (t *TestTaskA) GetSpec() string {
	return "*/5  *  *  *  *  *"
}

func (t *TestTaskA) GetDesc() string {
	return "测试任务A（每5秒执行一次）"
}

func (t *TestTaskA) Execute() {
	fmt.Println("------------测试任务A（每5秒执行一次）-------------")
}
