package core

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
)

func (mgr *Manager) GetBluePrintList() ([]*types.BluePrint, error) {
	var bluePrints []*types.BluePrint
	if err := mgr.storage.Find(&bluePrints, ""); err != nil {
		return []*types.BluePrint{}, errors.Trace(err)
	}
	return bluePrints, nil
}

func (mgr *Manager) GetBluePrintById(id int) (*types.BluePrint, error) {
	var bluePrint types.BluePrint
	if err := mgr.storage.FindOne(&bluePrint, "id=?", id); err != nil {
		return nil, errors.NotFoundf("blue print id %d", id)
	}
	return &bluePrint, nil
}
