package task

import (
	"beegostudy/service/cron"
	"fmt"
)

func init() {
	// cron.Task("Test", "*/2	*	*	*	*	*", Test)
	// cron.Task("Test2", "*/3	*	*	*	*	*", Test2)
	cron.Inter("asd", "qweqwe")
	fmt.Println("asd-------------------")
}

func Test() {
	fmt.Println("------------asdasd-------------")
}

func Test2() {
	fmt.Println("------------123123-------------")
}
