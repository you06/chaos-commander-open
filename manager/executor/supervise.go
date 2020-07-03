package executor

import (
	"github.com/juju/errors"
	"github.com/ngaut/log"
)

func (e *Executor) supervise() {
Supervise:
	for {
		select {
		case event := <-e.eventch:
			if !e.fatal {
				switch event {
				case EventPrepare:
					e.prepare()
				case EventAnsibleDownload:
					e.ansibleDownload()
				case EventAnsibleStop:
					e.ansibleStop()
				case EventAnsibleClearLog:
					e.ansibleClearLog()
				case EventAnsibleSwapOff:
					e.ansibleSwapOff()
				case EventAnsibleNTP:
					e.ansibleNTP()
				case EventAnsibleCPU:
					e.ansibleCPU()
				case EventAnsibleDeploy:
					e.ansibleDeploy()
				case EventAnsibleStart:
					e.ansibleStart()
				case EventLoadStart:
					e.loadStart()
				case EventLoadStop:
					e.loadStop()
				case EventMetricsStart:
					e.metricsWatcherStart()
				case EventMetricsStop:
					e.metricsWatcherStop()
				case EventBlueprintStart:
					e.blueprintStart()
				case EventBlueprintStop:
					e.blueprintStop()
				case EventFetchLog:
					log.Info("ready to fetch log")
					e.fetchAllLog()
				case EventStop:
					e.Stop()
				case EventStopSupervise:
					log.Info("stop supervise")
					break Supervise
				}
			}
		case err := <-e.errch:
			log.Errorf("error %v", errors.ErrorStack(err))
		}
	}
}
