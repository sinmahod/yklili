// This library implements a cron spec parser and runner.  See the README for
// more details.
package cron

import (
	"encoding/xml"
	"fmt"
	"log"
	"runtime"
	"sort"
	"time"
)

// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may be started, stopped, and the entries may
// be inspected while running.
type cron struct {
	XMLName  xml.Name       `xml:"TaskList"`
	Entries  []*Entry       `xml:"task"`
	stop     chan struct{}  `xml:"-"`
	add      chan *Entry    `xml:"-"`
	snapshot chan []*Entry  `xml:"-"`
	stopent  chan string    `xml:"-"` //暂停
	startent chan string    `xml:"-"` //恢复
	execute  chan string    `xml:"-"` //执行
	running  bool           `xml:"-"`
	ErrorLog *log.Logger    `xml:"-"`
	location *time.Location `xml:"-"`
}

// Job is an interface for submitted cron jobs.
type Job interface {
	Run()
}

// The Schedule describes a job's duty cycle.
type Schedule interface {
	// Return the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next(time.Time) time.Time
}

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
	//xml节点标签名称
	XMLName xml.Name `xml:"task"`

	//任务ID
	Id string `xml:"id,attr"`

	//任务状态
	Status bool `xml:"status,attr"`

	//描述
	Desc string `xml:",innerxml"`

	Type string `xml:"type,attr"`

	// The schedule on which this job should be run.
	schedule Schedule `xml:"-"`

	// The Job to run.
	job Job `xml:"-"`

	// The next time the job will run. This is the zero time if Cron has not been
	// started or this entry's schedule is unsatisfiable
	Next time.Time `xml:"-"`

	// The last time this job was run. This is the zero time if the job has never
	// been run.
	Prev time.Time `xml:"-"`
}

// byTime is a wrapper for sorting the entry array by time
// (with zero time at the end).
type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	// Two zero times should return false.
	// Otherwise, zero is "greater" than any other time.
	// (To sort it at the end of the list.)
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}

//按Id排序任务
type byId []*Entry

func (s byId) Len() int      { return len(s) }
func (s byId) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byId) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

// New returns a new Cron job runner, in the Local time zone.
func New() *cron {
	return NewWithLocation(time.Now().Location())
}

// NewWithLocation returns a new Cron job runner.
func NewWithLocation(location *time.Location) *cron {
	return &cron{
		Entries:  nil,
		add:      make(chan *Entry),
		stop:     make(chan struct{}),
		snapshot: make(chan []*Entry),
		stopent:  make(chan string),
		startent: make(chan string),
		execute:  make(chan string),
		running:  false,
		ErrorLog: nil,
		location: location,
	}
}

// A wrapper that turns a func() into a cron.Job
type FuncJob func()

func (f FuncJob) Run() { f() }

// AddFunc adds a func to the Cron to be run on the given schedule.
func (c *cron) addFunc(id, spec string, cmd func(), isrun bool) error {
	return c.addJob(id, "", "", spec, FuncJob(cmd), isrun)
}

// AddJob adds a Job to the Cron to be run on the given schedule.
func (c *cron) addJob(id, ttype, desc, spec string, cmd Job, isrun bool) error {
	schedule, err := Parse(spec)
	if err != nil {
		return err
	}
	return c.schedule(id, ttype, desc, schedule, cmd, isrun)
}

// Schedule adds a Job to the Cron to be run on the given schedule.
func (c *cron) schedule(id, ttype, desc string, schedule Schedule, cmd Job, isrun bool) error {
	//检查ID是否存在
	for _, e := range c.Entries {
		if id == e.Id {
			return fmt.Errorf("这个任务已经存在: %s", id)
		}
	}

	entry := &Entry{
		schedule: schedule,
		job:      cmd,
		Id:       id,
		Status:   isrun,
		Type:     ttype,
		Desc:     desc,
	}
	if !c.running {
		c.Entries = append(c.Entries, entry)
		return nil
	}

	c.add <- entry
	return nil
}

// 暂停某个定时任务
func (c *cron) stopFunc(id string) {
	if !c.running {
		for _, e := range c.Entries {
			if id == e.Id {
				e.Status = false
				return
			}
		}
	} else {
		c.stopent <- id
		<-c.stopent
	}
}

// 启动某个定时任务
func (c *cron) startFunc(id string) {
	if !c.running {
		for _, e := range c.Entries {
			if id == e.Id {
				e.Status = true
				return
			}
		}
	}

	c.startent <- id
	<-c.startent
}

// 执行某个定时任务
func (c *cron) execFunc(id string) error {
	if !c.running {
		return fmt.Errorf("主线程已经停止无法执行任务")
	}
	c.execute <- id
	<-c.execute
	return nil
}

// 查询某个定时任务当前运行状态
func (c *cron) funcStatus(id string) bool {
	if !c.running {
		return false
	}
	c.snapshot <- nil
	entries := <-c.snapshot
	for _, e := range entries {
		if e.Id == id {
			return e.Status
		}
	}
	return false
}

// 任务列表
func (c *cron) EntryList() []*Entry {
	if c.running {
		c.snapshot <- nil
		x := <-c.snapshot
		return x
	}
	return c.entrySnapshot()
}

// Location gets the time zone location
func (c *cron) Location() *time.Location {
	return c.location
}

// 启动线程
func (c *cron) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

// 总线程状态
func (c *cron) status() bool {
	return c.running
}

func (c *cron) runWithRecovery(j Job) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.logf("cron: panic running job: %v\n%s", r, buf)
		}
	}()
	j.Run()
}

// Run the scheduler.. this is private just due to the need to synchronize
// access to the 'running' state variable.
func (c *cron) run() {
	// Figure out the next activation times for each entry.
	now := time.Now().In(c.location)
	for _, entry := range c.Entries {
		if !entry.Status {
			t, _ := time.Parse("2006-01-02 15:04:05", "2222-02-22 22:22:22")
			entry.Next = entry.schedule.Next(t)
		} else {
			entry.Next = entry.schedule.Next(now)
		}
	}

	for {
		// Determine the next entry to run.
		sort.Sort(byTime(c.Entries))

		var effective time.Time
		if len(c.Entries) == 0 || c.Entries[0].Next.IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = c.Entries[0].Next
		}

		timer := time.NewTimer(effective.Sub(now))
		select {
		case now = <-timer.C:
			now = now.In(c.location)
			// Run every entry whose next time was this effective time.
			for _, e := range c.Entries {
				if e.Next != effective || !e.Status {
					break
				}
				log.Printf("开始执行任务:" + e.Desc)
				go c.runWithRecovery(e.job)
				e.Prev = e.Next
				e.Next = e.schedule.Next(now)
			}
			continue

		case newEntry := <-c.add:
			c.Entries = append(c.Entries, newEntry)
			if newEntry.Status {
				t, _ := time.Parse("2006-01-02 15:04:05", "2222-02-22 22:22:22")
				newEntry.Next = newEntry.schedule.Next(t)
			} else {
				newEntry.Next = newEntry.schedule.Next(time.Now().In(c.location))
			}

		//暂停任务
		case stopEntry := <-c.stopent:
			for _, e := range c.Entries {
				if e.Id == stopEntry {
					e.Status = false
					t, _ := time.Parse("2006-01-02 15:04:05", "2222-02-22 22:22:22")
					e.Next = e.schedule.Next(t)
					go c.saveTaskFile()
					break
				}
			}
			c.stopent <- ""

		//恢复任务
		case startEntry := <-c.startent:
			for _, e := range c.Entries {
				if e.Id == startEntry {
					e.Status = true
					e.Next = e.schedule.Next(time.Now().In(c.location))
					go c.saveTaskFile()
					break
				}
			}
			c.startent <- ""

		//立即执行
		case exeEntryId := <-c.execute:
			for _, e := range c.Entries {
				if e.Id == exeEntryId {
					go c.runWithRecovery(e.job)
					e.Prev = time.Now()
					break
				}
			}
			c.execute <- ""

		case <-c.snapshot:
			c.snapshot <- c.entrySnapshot()

		case <-c.stop:
			timer.Stop()
			return
		}

		// 'now' should be updated after newEntry and snapshot cases.
		now = time.Now().In(c.location)
		timer.Stop()
	}
}

// Logs an error to stderr or to the configured error log
func (c *cron) logf(format string, args ...interface{}) {
	if c.ErrorLog != nil {
		c.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// 停止线程
func (c *cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

// entrySnapshot returns a copy of the current cron entry list.
func (c *cron) entrySnapshot() []*Entry {
	entries := []*Entry{}
	for _, e := range c.Entries {
		entries = append(entries, &Entry{
			XMLName:  e.XMLName,
			schedule: e.schedule,
			Next:     e.Next,
			Prev:     e.Prev,
			job:      e.job,
			Id:       e.Id,
			Status:   e.Status,
			Type:     e.Type,
			Desc:     e.Desc,
		})
	}
	return entries
}
