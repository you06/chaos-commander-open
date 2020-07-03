package executor

import (
	"context"
	"fmt"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/manager/blueprint"
	"github.com/you06/chaos-commander/manager/core"
	"github.com/you06/chaos-commander/manager/executor/content"
	"github.com/you06/chaos-commander/manager/metrics"
	logger "github.com/you06/chaos-commander/pkg/log"
	"github.com/you06/chaos-commander/pkg/types"
	"github.com/you06/chaos-commander/util"
	"os"
	"runtime/debug"
	"time"
)

const WAIT_FOR_LOAD = time.Duration(10)

type contextCancel struct {
	ctx    context.Context
	cancel context.CancelFunc
	finish bool
}

type Executor struct {
	mgr           *core.Manager
	content       *content.Content
	history       *types.History
	logPath       string
	logger        *logger.Log
	status        ExecStatus
	loadCtx       *contextCancel
	blueprintCtx  context.Context
	blueprintExec *blueprint.Blueprint
	metrics       *metrics.Metrics
	eventch       chan EventType
	errch         chan error
	fatal         bool
	callback      func()
}

func New(mgr *core.Manager, job *types.Job, callback func()) *Executor {
	loadCtx, loadCancel := context.WithCancel(context.Background())
	return &Executor{
		mgr: mgr,
		content: &content.Content{
			Pkg: mgr.Pkg,
			Job: job,
		},
		status: StatusPending,
		loadCtx: &contextCancel{
			ctx: loadCtx,
			cancel: loadCancel,
			finish: false,
		},
		eventch: make(chan EventType, 1),
		errch: make(chan error),
		fatal: false,
		callback: callback,
	}
}

func (e *Executor)Start() {
	if err := e.prepare(); err != nil {
		e.fatalError(errors.Trace(err))
		return
	}

	if !e.content.Job.GetSkipDeploy() {
		// Download latest version
		if err := e.ansibleDownload(); err != nil {
			e.fatalError(errors.Trace(err))
			return
		}
		// Stop cluster
		if err := e.ansibleStop(); err != nil {
			e.logger.Warnf("%v", errors.ErrorStack(err))
		}
		// Clear log
		if err := e.ansibleClearLog(); err != nil {
			e.fatalError(errors.Trace(err))
			return
		}
		// Adjust CPU and make sure ntp works
		if err := e.ansibleSwapOff(); err != nil {
			e.fatalError(errors.Trace(err))
			return
		}
		if err := e.ansibleNTP(); err != nil {
			e.fatalError(errors.Trace(err))
			return
		}
		if err := e.ansibleCPU(); err != nil {
			e.fatalError(errors.Trace(err))
			return
		}
		// Run cluster
		if err := e.ansibleDeploy(); err != nil {
			e.fatalError(errors.Trace(err))
			return
		}
		if err := e.ansibleStart(); err != nil {
			e.fatalError(errors.Trace(err))
			return
		}
	}
	// Start metrics check
	//e.metricsWatcherStart()
	// Run load
	e.loadStart()
	if e.content.Job.GetPessimistic() {
		go e.switchPessimistic()
	}
	time.Sleep(WAIT_FOR_LOAD * time.Second)
	// Run blue print
	if err := e.blueprintStart(); err != nil {
		e.fatalError(errors.Trace(err))
		return
	}
}

func (e *Executor)Stop() {
	log.Info("stop execution")
	e.loadStop()
	log.Info("before fetch log.")
	e.fetchAllLog()
	//e.metricsWatcherStop()
	//e.eventch <- EventAnsibleStop

	if err := e.makeResult(); err != nil {
		log.Error(err)
	}

	//e.eventch <- EventStopSupervise
	e.callback()
}

func (e *Executor)fatalError(err error) {
	e.fatal = true
	log.Errorf("fatal error %v", errors.ErrorStack(err))
	debug.PrintStack()
	e.Stop()
}

func (e *Executor)prepare() error {
	// get ansible
	ansible, err := e.mgr.GetAnsibleById(e.content.Job.GetAnsibleID())
	if err != nil {
		return errors.Trace(errors.Trace(err))
	}
	e.content.Ansible = ansible
	e.metrics = metrics.New(ansible.GetProm())

	// get load
	load, err := e.mgr.GetLoadById(e.content.Job.GetLoadID())
	if err != nil {
		return errors.Trace(errors.Trace(err))
	}
	e.content.Load = load

	// get blueprint
	blue, err := e.mgr.GetBluePrintById(e.content.Job.GetBluePrintID())
	if err != nil {
		return errors.Trace(errors.Trace(err))
	}
	e.content.Blueprint = blue

	// get resources
	resourceIds := e.content.Job.GetResources()
	var resources []*types.Resource
	for _, resourceId := range resourceIds {
		resource, err := e.mgr.GetResourceById(resourceId)
		if err != nil {
			return errors.Trace(errors.Trace(err))
		} else {
			resources = append(resources, resource)
		}
	}

	// prepare history
	if err := e.createResult(); err != nil {
		return errors.Trace(errors.Trace(err))
	}

	// prepare log path
	err = e.prepareLogPath()
	if err != nil {
		return errors.Trace(errors.Trace(err))
	}

	e.content.Logger = logger.New(e.content.Pkg.LogEngine, logger.TypeManager, e.content.Job.GetID(), e.history.GetID(), "manager")
	e.logger = e.content.Logger
	e.content.Resources = resources
	e.status = StatusPrepare
	return nil
}

func (e *Executor)prepareLogPath() error {
	logPath := fmt.Sprintf("%s/%s-%s", e.mgr.Config.LogPath, e.content.Job.GetName(), util.Now())
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		return errors.Trace(err)
	}
	if err := os.MkdirAll(logPath, 0777); err != nil {
		return errors.Trace(err)
	}
	e.logPath = logPath
	//e.logger = logger.New(e.logPath)
	e.history.Log = logPath
	return nil
}
