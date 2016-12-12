package progress

import "fmt"

var threadtask = make(map[string]*ProgressTask)

type ProgressTask struct {
	id   string
	perc int
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
		return fmt.Errorf("任务[%s]已经在执行中，请等待完成，当前进度为(%d%s)", t.id, tt.perc, "%")
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
	return t.perc
}

// 得到任务当前进度
func (t *ProgressTask) SetPerc(perc int) {
	t.perc = perc
}
