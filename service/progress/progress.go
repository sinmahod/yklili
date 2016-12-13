package progress

import "fmt"

var threadtask = make(map[string]*ProgressTask)

func GetPerc(taskId string) *ProgressTask {
	if t, ok := threadtask[taskId]; ok {
		return t
	} else {
		return nil
	}
}

type ProgressTask struct {
	id   string
	Perc int
	Msg  string
	fc   func()
}

// 设置任务ID
func (t *ProgressTask) SetTaskId(taskId string) error {
	t.id = taskId
	return nil
}

func (t *ProgressTask) Start() error {
	if t.id == "" {

	}
	if tt, ok := threadtask[t.id]; ok {
		return fmt.Errorf("任务[%s]已经在执行中，请等待完成，当前进度为(%d%s)", t.id, tt.Perc, "%")
	}
	go t.execute()
	threadtask[t.id] = t
	return nil
}

func (t *ProgressTask) execute() {
	t.fc()
	delete(threadtask, t.id)
}

func (t *ProgressTask) SetFunc(fc func()) {
	t.fc = fc
}

// 得到任务当前进度
func (t *ProgressTask) GetPerc() int {
	return t.Perc
}

// 得到任务当前进度
func (t *ProgressTask) SetPerc(perc int) {
	t.Perc = perc
}

// 得到任务当前消息
func (t *ProgressTask) GetMsg() string {
	return t.Msg
}

// 得到任务当前消息
func (t *ProgressTask) SetMsg(format string, a ...interface{}) {
	t.Msg = fmt.Sprintf(format, a...)
}
