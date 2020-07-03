package core

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
)

func (mgr *Manager) GetAnsibleList() ([]*types.Ansible, error) {
	var ansibles []*types.Ansible
	if err := mgr.storage.Find(&ansibles, ""); err != nil {
		return []*types.Ansible{}, errors.Trace(err)
	}
	return ansibles, nil
}

func (mgr *Manager) GetAnsibleById(id int) (*types.Ansible, error) {
	var ansible types.Ansible
	if err := mgr.storage.FindOne(&ansible, "id=?", id); err != nil {
		return nil, errors.NotFoundf("ansible id %d", id)
	}
	return &ansible, nil
}
