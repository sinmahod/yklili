/**
 *   Github robfig/cron 封装
 *   提供定时任务常用方法
 */
package cron

import (
	"beegostudy/util"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

var crontab *Cron

func init() {
	crontab = New()
}

var sss = make(map[string]string)

func Inter(id, value string) {
	sss[id] = value
}

/**
 *   新创建任务
 *   id 		任务ID
 *   spec 	crontab表达式 ： 3  *  *  *  *  *
 *   fc		函数值
 *   desc 	任务描述
 */
func Task(id, spec string, fc func(), desc ...string) {
	crontab.AddFunc(id, spec, fc)
}

/**
 *   启动线程
 */
func RunCron() {
	crontab.Start()
}

/**
 *   停止线程
 */
func StopCron() {
	crontab.Stop()
}

/**
 *   线程运行状态
 */
func CronStatus() bool {
	return crontab.Status()
}

/**
 *   任务状态
 *   id	任务id
 */
func TaskStatus(id string) bool {
	return crontab.FuncStatus(id)
}

/**
 *   运行任务
 */
func TaskStart(id string) {
	crontab.StartFunc(id)
}

/**
 *   停止任务
 */
func TaskStop(id string) {
	crontab.StopFunc(id)
}

/**
 *   执行任务
 */
func TaskExecute(id string) {
	crontab.ExecFunc(id)
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
	//读
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
