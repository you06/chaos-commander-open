package core

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
)

func (mgr *Manager) GetJobsList() ([]*types.Job, error) {
	var jobs []*types.Job
	if err := mgr.storage.Find(&jobs, ""); err != nil {
		return []*types.Job{}, errors.Trace(err)
	}
	return jobs, nil
}

func (mgr *Manager) GetJobById(id int) (*types.Job, error) {
	var job types.Job
	if err := mgr.storage.FindOne(&job, "id=?", id); err != nil {
		return nil, errors.NotFoundf("job id %d", id)
	}
	return &job, nil
}
