package executor

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/manager/blueprint"
	"time"
)

func (e *Executor)blueprintStart() error {
	blueprintExec, err := blueprint.New(e.content, e.content.Resources, func() {
		time.Sleep(5 * time.Minute)
		//e.eventch <- EventBlueprintStop
		e.blueprintStop()
	})
	if err != nil {
		return errors.Trace(err)
	}
	e.blueprintExec = blueprintExec
	e.blueprintExec.Start()
	return nil
}

func (e *Executor)blueprintStop() {
	//e.eventch <- EventLoadStop
	e.Stop()
	//e.eventch <- EventStop
}
