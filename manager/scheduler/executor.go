package scheduler

import (
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/manager/executor"
	"github.com/you06/chaos-commander/pkg/types"
)

func (sch *Scheduler) Exec(jobID int) {
	job, err := sch.mgr.GetJobById(jobID)
	if err != nil {
		log.Error(errors.ErrorStack(err))
		return
	}
	log.Infof("job ID: %d, job name: %s ready to execute", jobID, job.GetName())
	// wait for resource is ready
	ready := false
	for !ready {
		if resourceReady, err := sch.wait(job); err != nil {
			log.Error(errors.ErrorStack(err))
			return
		} else if resourceReady {
			ready = true
		} else {
			time.Sleep(time.Minute)
		}
	}

	log.Infof("job %d start", job.GetID())
	exec := executor.New(sch.mgr, job, func() {
		log.Infof("job %d finished", job.GetID())
		if err := sch.clearLock(job); err != nil {
			log.Error(err)
		}
	})
	go exec.Start()
}

func (sch *Scheduler) wait(job *types.Job) (bool, error) {
	sch.Lock()
	defer sch.Unlock()
	resourceIDs := job.GetResources()

	resources, err := sch.mgr.GetResourceByIds(resourceIDs)
	if err != nil {
		return false, errors.Trace(err)
	}

	for _, resource := range resources {
		if resource.GetStatus() == types.ResourceStatusWait {
			return false, nil
		}
	}

	if err := sch.mgr.LockResource(resourceIDs); err != nil {
		return false, errors.Trace(err)
	}

	return true, nil
}

func (sch *Scheduler) clearLock(job *types.Job) error {
	if err := sch.mgr.FreeResource(job.GetResources()); err != nil {
		return errors.Trace(err)
	}
	return nil
}
