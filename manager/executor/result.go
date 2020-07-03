package executor

import (
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
)

func (e *Executor) makeResult() error {
	if err := e.saveResult(); err != nil {
		return errors.Trace(err)
	}
	if err := e.slackNotice(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *Executor) createResult() error {
	e.history = &types.History{
		JobID:      e.content.Job.GetID(),
		Status:     types.HistoryStatusPending,
		FinishedAt: time.Now(),
	}
	if err := e.mgr.CreateHistory(e.history); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *Executor) saveResult() error {
	var status types.HistoryStatus
	if e.fatal {
		status = types.HistoryStatusFailed
	} else {
		status = types.HistoryStatusSuccess
	}
	e.history.Status = status
	e.history.FinishedAt = time.Now()
	if err := e.mgr.UpdateHistory(e.history); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// TODO: provided log URL in slack message
func (e *Executor) slackNotice() error {
	var msg string

	if !e.fatal {
		msg = fmt.Sprintf("job %d pass", e.content.Job.GetID())
	} else {
		msg = fmt.Sprintf("job %d failed", e.content.Job.GetID())
	}

	if err := e.mgr.Pkg.Slack.SendMessage(msg); err != nil {
		return errors.Trace(err)
	}

	return nil
}
