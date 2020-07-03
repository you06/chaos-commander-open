package executor

import (
	"context"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/pkg/mysql"
)

func (e *Executor) switchPessimistic() {
	status := false
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-ticker.C:
			e.doSwitchPessimistic(!status)
			status = !status
		}
	}
}

func (e *Executor) doSwitchPessimistic(status bool) {
	e.loadCtx.finish = true
	e.loadCtx.cancel()
	// switch pessimistic
	db, err := mysql.OpenDB(e.content.Job.GetDSN(), 1)
	if err != nil {
		e.logger.Errorf("connect db error %v", errors.ErrorStack(err))
		log.Errorf("connect db error %v", errors.ErrorStack(err))
	}
	if status {
		e.logger.Infof("turn on pessimistic %s", e.content.Job.GetDSN())
		log.Infof("turn on pessimistic %s", e.content.Job.GetDSN())
		db.MustExec("set @@global.tidb_txn_mode = 'pessimistic'")
	} else {
		e.logger.Infof("turn off pessimistic %s", e.content.Job.GetDSN())
		log.Infof("turn off pessimistic %s", e.content.Job.GetDSN())
		db.MustExec("set @@global.tidb_txn_mode = 'optimistic'")
	}
	if err := db.CloseDB(); err != nil {
		e.logger.Errorf("close db error %v", err)
		log.Errorf("close db error %v", err)
	}
	// sleep for waiting load quit and pessimistic switch done
	time.Sleep(120 * time.Second)
	// prepare a new one
	loadCtx, loadCancel := context.WithCancel(context.Background())
	e.loadCtx = &contextCancel{
		ctx:    loadCtx,
		cancel: loadCancel,
		finish: false,
	}
	e.loadStart()
}
