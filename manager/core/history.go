package core

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
)

func (mgr *Manager) GetHistoryList() ([]*types.History, error) {
	var histories []*types.History
	if err := mgr.storage.Find(&histories, ""); err != nil {
		return []*types.History{}, errors.Trace(err)
	}
	return histories, nil
}

func (mgr *Manager) GetHistoryById(id int) (*types.History, error) {
	var history types.History
	if err := mgr.storage.FindOne(&history, "id=?", id); err != nil {
		return nil, errors.NotFoundf("history id %d", id)
	}
	return &history, nil
}

// CreateHistory create a history
func (mgr *Manager) CreateHistory(history *types.History) error {
	if err := mgr.storage.Save(history); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// UpdateHistory update single history
func (mgr *Manager) UpdateHistory(history *types.History) error {
	return errors.Trace(mgr.storage.Save(history))
}
