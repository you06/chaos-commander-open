package scheduler

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/you06/chaos-commander/pkg/types"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	cron "github.com/robfig/cron/v3"
	"github.com/you06/chaos-commander/manager/core"
)

type Scheduler struct {
	sync.Mutex
	mgr  *core.Manager
	cron *cron.Cron
}

// New creates the scheduler
func New(mgr *core.Manager) *Scheduler {
	sch := Scheduler{
		mgr:  mgr,
		cron: cron.New(),
	}

	if err := sch.Prepare(); err != nil {
		log.Errorf("scheduler prepare failed %v", err)
	}
	return &sch
}

// Prepare init the cron
func (sch *Scheduler) Prepare() error {
	if err := sch.mgr.ClearResourceLock(); err != nil {
		return errors.Trace(err)
	}
	if err := sch.Schedule(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Start the cron scheduler
func (sch *Scheduler) Start() {
	log.Info("start scheduler")
	sch.cron.Start()
}

func (sch *Scheduler) Stop() {
	sch.cron.Stop()
}

// Schedule make schedule plan from database
func (sch *Scheduler) Schedule() error {
	jobs, err := sch.mgr.GetJobsList()
	if err != nil {
		return errors.Trace(err)
	}
	for index, job := range jobs {
		go func(index int, job *types.Job) {
			cronPlan := job.GetRunningPlan()
			if cronPlan == "never" {
				return
			}
			switch cronPlan {
			case "hourly", "daily", "midnight", "weekly", "monthly", "yearly", "annually":
				cronPlan = fmt.Sprintf("@%s", cronPlan)
			}

			jobID := job.GetID()

			delay := time.Duration(index*15) * time.Minute
			log.Infof("%s, job ID: %d, job name: %s, delay: %s", cronPlan, job.GetID(), job.GetName(), delay)
			time.Sleep(delay)

			sch.cron.AddFunc(cronPlan, func() {
				sch.Exec(jobID)
			})
			if strings.HasPrefix(cronPlan, "@every") {
				go sch.Exec(jobID)
			}
		}(index, job)
	}
	return nil
}

// ReSchedule clear schedule plan and make a new one from database
func (sch *Scheduler) ReSchedule() error {
	sch.clear()
	if err := sch.Schedule(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (sch *Scheduler) clear() {
	entries := sch.cron.Entries()
	for _, entry := range entries {
		sch.cron.Remove(entry.ID)
	}
}

// Test run single job immediately
func (sch *Scheduler) Test(id int) error {
	job, err := sch.mgr.GetJobById(id)
	if err != nil {
		return err
	}
	sch.Exec(job.GetID())
	return nil
}
