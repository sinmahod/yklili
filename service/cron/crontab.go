/**
 *   Github robfig/cron 封装
 *   提供定时任务常用方法
 */
package cron

import (
	"encoding/xml"
	"fmt"
	"github.com/sinmahod/yklili/util/fileutil"
	"sort"
	"strings"

	"github.com/astaxie/beego"
)

var crontab *cron

func init() {
	crontab = New()
	readTaskFile()
}

/**
 *   任务接口
 */
type TaskInterface interface {
	// 任务ID，唯一
	GetId() string
	// crontab 表达式(秒 分 时 日 月 周)
	GetSpec() string
	// 任务描述
	GetDesc() string
	// 定时执行的任务
	Execute()
}

var taskstatus map[string]bool

/**
 *   注册任务
 *   task 		任务接口实现
 */
func RegisterTask(task TaskInterface) {
	if status, ok := taskstatus[task.GetId()]; ok {
		if status {
			crontab.addTask(task, true)
		} else {
			crontab.addTask(task, false)
		}
	} else {
		crontab.addTask(task, true)
	}
}

/**
 *   启动线程
 */
func RunCron() {
	crontab.saveTaskFile()
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
	return crontab.status()
}

/**
 *   任务状态
 *   id	任务id
 */
func TaskStatus(id string) bool {
	return crontab.funcStatus(id)
}

/**
 *   运行任务
 */
func TaskStart(id string) {
	crontab.startFunc(id)
}

/**
 *   停止任务
 */
func TaskStop(id string) {
	crontab.stopFunc(id)
}

/**
 *   执行任务
 */
func TaskExecute(ids string) {
	idx := strings.Split(ids, ",")
	for _, id := range idx {
		crontab.execFunc(id)
	}
}

/**
 *   获取所有任务列表
 */
func GetTaskList() []*Entry {
	entries := crontab.EntryList()
	sort.Sort(byId(entries))
	return entries
}

var taskFile = beego.AppPath + "/conf/task.xml"

// 读取定时任务配置文件
func readTaskFile() {
	taskstatus = make(map[string]bool)
	var t *cron
	if fileutil.Exist(taskFile) {
		fileutil.XMLToStruct(taskFile, &t)
		if t != nil {
			for _, entry := range t.Entries {
				taskstatus[entry.Id] = entry.Status
			}
		}

	}
}

// 添加任务
func (c *cron) addTask(task TaskInterface, run bool) {
	c.addJob(task.GetId(), fmt.Sprintf("%T", task), task.GetDesc(), task.GetSpec(), FuncJob(task.Execute), run)
}

// 保存到文件
func (c *cron) saveTaskFile() {
	t := &cron{}
	t.XMLName = xml.Name{"", "TaskList"}
	t.Entries = c.EntryList() //得到最新的
	fileutil.XMLStructToFile(taskFile, t)
}
