package core

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
)

func (mgr *Manager) GetLoadList() ([]*types.Load, error) {
	var loads []*types.Load
	if err := mgr.storage.Find(&loads, ""); err != nil {
		return []*types.Load{}, errors.Trace(err)
	}
	return loads, nil
}

func (mgr *Manager) GetLoadById(id int) (*types.Load, error) {
	var load types.Load
	if err := mgr.storage.FindOne(&load, "id=?", id); err != nil {
		return nil, errors.NotFoundf("load id %d", id)
	}
	return &load, nil
}
