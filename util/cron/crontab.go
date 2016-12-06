package cron

import (
	"beegostudy/util"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

type tasks map[string]timetask

type timetask struct {
	desc   string
	cronEx string //表达式  分时日月年周
	*Cron
}

var taskList = make(tasks)

/**
 *  创建任务,返回任务列表
 */
func Task(id, cronEx string, fc func(), desc string) error {
	task := new(timetask)
	c := New()
	c.AddFunc(cronEx, fc)
	task.desc = desc
	task.cronEx = cronEx
	task.Cron = c
	if _, ok := taskList[id]; !ok {
		taskList[id] = *task
	} else {
		return fmt.Errorf("任务%s已经存在", id)
	}
	return nil
}

type TaskXml struct {
	XMLName xml.Name   `xml:"tasks"`
	Node    []TaskNode `xml:"task"`
}

type TaskNode struct {
	XMLName xml.Name `xml:"task"`
	Id      string   `xml:"id,attr"`
	Status  bool     `xml:",innerxml"`
}

func ReadTaskFile() {
	var result TaskXml
	err := util.XMLToStruct(beego.AppPath+"/conf/task.xml", &result)
	fmt.Println(err)
	fmt.Println(result)
	fmt.Println(result.Node)
	for _, o := range result.Node {
		fmt.Println(o.Id+"===", o.Status)
	}

	//写
	for i, line := range result.Node {
		//修改ApplicationName节点的内部文本innerText
		if strings.EqualFold(line.Id, "ApplicationName") {
			//注意修改的不是line对象，而是直接使用result中的真实对象
			result.Node[i].Status = false
		}
	}
	util.XMLStructToFile(beego.AppPath+"/conf/task.xml", &result)
}

func createTaskXML() {

}

/**
 *  返回任务列表
 */
func GetTasks() tasks {
	return taskList
}

/**
 *  启动所有任务
 */
func StartTasks() {
	for _, task := range taskList {
		task.Cron.Start()
	}
}

/**
 *  停止所有任务
 */
func StopTasks() {
	for _, task := range taskList {
		task.Cron.Stop()
	}
}

/**
 *  停止某个任务
 */
func StopTask(id string) error {
	if task, ok := taskList[id]; ok {
		if task.Cron.Status() {
			task.Cron.Stop()
		}
	} else {
		return fmt.Errorf("任务%s不存在", id)
	}
	return nil
}

/**
 *  启动某个任务
 */
func StartTask(id string) error {
	if task, ok := taskList[id]; ok {
		if !task.Cron.Status() {
			task.Cron.Start()
		}
	} else {
		return fmt.Errorf("任务%s不存在", id)
	}
	return nil
}
