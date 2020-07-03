package executor

import (
	"fmt"
	"github.com/juju/errors"
	//_log "github.com/ngaut/log"
	"github.com/you06/chaos-commander/pkg/log"
	"github.com/you06/chaos-commander/pkg/types"
	"github.com/you06/chaos-commander/util"
	//"os/exec"
	"path"
	"strings"
	//"syscall"
	"time"
)

func (e *Executor)loadStart() {
	logPath := path.Join(e.logPath, fmt.Sprintf("%s-load_%s.log", e.content.Job.GetName(), time.Now().Format("2006-01-02T15:04:05")))
	fmt.Println("load log path", e.mgr.Config.LogPath, logPath)
	switch e.content.Load.GetType() {
	case types.LoadTypePath:
		go func () {
			logger := log.New(e.content.Pkg.LogEngine, log.TypeLoad, e.content.Job.GetID(), e.history.GetID(), "load")
			err := util.DoCmdContextWithLogger(e.loadCtx.ctx, logger, e.mgr.Config.LogPath,
				e.content.Load.GetPath(), strings.Split(e.content.Load.GetParam(), " ")...)
			logger.Quit()

			if !e.loadCtx.finish {
				e.fatalError(errors.Trace(err))
			}
			//if execExitError, ok := err.(*exec.ExitError); ok {
			//	if status, ok := execExitError.Sys().(syscall.WaitStatus); ok {
			//		if status != 9 {
			//			_log.Error(errors.ErrorStack(err))
			//			e.fatalError(errors.Trace(err))
			//		}
			//	}
			//} else if err != nil {
			//	_log.Error(errors.ErrorStack(err))
			//	e.fatalError(errors.Trace(err))
			//}
		}()
	}
}

func (e *Executor)loadStop() {
	e.loadCtx.finish = true
	e.loadCtx.cancel()
}
