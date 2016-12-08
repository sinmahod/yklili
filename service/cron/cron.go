// This library implements a cron spec parser and runner.  See the README for
// more details.
package cron

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"time"
)

// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may be started, stopped, and the entries may
// be inspected while running.
type Cron struct {
	entries  []*Entry
	stop     chan struct{}
	add      chan *Entry
	snapshot chan []*Entry
	stopent  chan string //暂停
	startent chan string //恢复
	execute  chan string //执行
	running  bool
	ErrorLog *log.Logger
	location *time.Location
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
	// The schedule on which this job should be run.
	Schedule Schedule

	// The next time the job will run. This is the zero time if Cron has not been
	// started or this entry's schedule is unsatisfiable
	Next time.Time

	// The last time this job was run. This is the zero time if the job has never
	// been run.
	Prev time.Time

	// The Job to run.
	Job Job

	//TaskId
	Id string

	//Status  true = runing
	Status bool
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

// New returns a new Cron job runner, in the Local time zone.
func New() *Cron {
	return NewWithLocation(time.Now().Location())
}

// NewWithLocation returns a new Cron job runner.
func NewWithLocation(location *time.Location) *Cron {
	return &Cron{
		entries:  nil,
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
func (c *Cron) AddFunc(id, spec string, cmd func()) error {
	return c.AddJob(id, spec, FuncJob(cmd))
}

// AddJob adds a Job to the Cron to be run on the given schedule.
func (c *Cron) AddJob(id, spec string, cmd Job) error {
	schedule, err := Parse(spec)
	if err != nil {
		return err
	}
	return c.Schedule(id, schedule, cmd)
}

// Schedule adds a Job to the Cron to be run on the given schedule.
func (c *Cron) Schedule(id string, schedule Schedule, cmd Job) error {
	//检查ID是否存在
	for _, e := range c.entries {
		if id == e.Id {
			return fmt.Errorf("这个任务已经存在: %s", id)
		}
	}
	entry := &Entry{
		Schedule: schedule,
		Job:      cmd,
		Id:       id,
		Status:   true,
	}
	if !c.running {
		c.entries = append(c.entries, entry)
		return nil
	}

	c.add <- entry
	return nil
}

// 暂停某个定时任务
func (c *Cron) StopFunc(id string) {
	if !c.running {
		for _, e := range c.entries {
			if id == e.Id {
				e.Status = false
				return
			}
		}
	} else {
		c.stopent <- id
	}
}

// 启动某个定时任务
func (c *Cron) StartFunc(id string) {
	if !c.running {
		for _, e := range c.entries {
			if id == e.Id {
				e.Status = true
				return
			}
		}
	}

	c.startent <- id
}

// 执行某个定时任务
func (c *Cron) ExecFunc(id string) error {
	if !c.running {
		return fmt.Errorf("主线程已经停止无法执行任务")
	}
	c.execute <- id
	return nil
}

// 查询某个定时任务当前运行状态
func (c *Cron) FuncStatus(id string) bool {
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

// 不需要对外开放
// Entries returns a snapshot of the cron entries.
// func (c *Cron) Entries() []*Entry {
// 	if c.running {
// 		c.snapshot <- nil
// 		x := <-c.snapshot
// 		return x
// 	}
// 	return c.entrySnapshot()
// }

// Location gets the time zone location
func (c *Cron) Location() *time.Location {
	return c.location
}

// 启动线程
func (c *Cron) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

// 总线程状态
func (c *Cron) Status() bool {
	return c.running
}

func (c *Cron) runWithRecovery(j Job) {
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
func (c *Cron) run() {
	// Figure out the next activation times for each entry.
	now := time.Now().In(c.location)
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		// Determine the next entry to run.
		sort.Sort(byTime(c.entries))

		var effective time.Time
		if len(c.entries) == 0 || c.entries[0].Next.IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = c.entries[0].Next
		}

		timer := time.NewTimer(effective.Sub(now))
		select {
		case now = <-timer.C:
			now = now.In(c.location)
			// Run every entry whose next time was this effective time.
			for _, e := range c.entries {
				if e.Next != effective || !e.Status {
					break
				}
				go c.runWithRecovery(e.Job)
				e.Prev = e.Next
				e.Next = e.Schedule.Next(now)
			}
			continue

		case newEntry := <-c.add:
			fmt.Println("Add")
			c.entries = append(c.entries, newEntry)
			newEntry.Next = newEntry.Schedule.Next(time.Now().In(c.location))

		//暂停任务
		case stopEntry := <-c.stopent:
			fmt.Println(stopEntry)
			for _, e := range c.entries {
				if e.Id == stopEntry {
					e.Status = false
					t, _ := time.Parse("2006-01-02 15:04:05", "2222-02-22 22:22:22")
					e.Next = e.Schedule.Next(t)
					break
				}
			}

		//恢复任务
		case startEntry := <-c.startent:
			for _, e := range c.entries {
				if e.Id == startEntry {
					e.Status = true
					e.Next = e.Schedule.Next(time.Now().In(c.location))
					break
				}
			}

		//立即执行
		case exeEntryId := <-c.execute:
			for _, e := range c.entries {
				if e.Id == exeEntryId {
					go c.runWithRecovery(e.Job)
					e.Prev = effective
					break
				}
			}

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
func (c *Cron) logf(format string, args ...interface{}) {
	if c.ErrorLog != nil {
		c.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// 停止线程
func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

// entrySnapshot returns a copy of the current cron entry list.
func (c *Cron) entrySnapshot() []*Entry {
	entries := []*Entry{}
	for _, e := range c.entries {
		entries = append(entries, &Entry{
			Schedule: e.Schedule,
			Next:     e.Next,
			Prev:     e.Prev,
			Job:      e.Job,
			Id:       e.Id,
			Status:   e.Status,
		})
	}
	return entries
}
