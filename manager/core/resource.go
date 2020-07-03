package core

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
)

// GetResourceList get resoruces list
func (mgr *Manager) GetResourceList() ([]*types.Resource, error) {
	var resources []*types.Resource
	if err := mgr.storage.Find(&resources, ""); err != nil {
		return []*types.Resource{}, errors.Trace(err)
	}
	return resources, nil
}

// GetResourceById get resource by id
func (mgr *Manager) GetResourceById(id int) (*types.Resource, error) {
	var resource types.Resource
	if err := mgr.storage.FindOne(&resource, "id=?", id); err != nil {
		return nil, errors.NotFoundf("resource id %d", id)
	}
	return &resource, nil
}

// GetResourceByIds get range resources
func (mgr *Manager) GetResourceByIds(ids []int) ([]*types.Resource, error) {
	var resources []*types.Resource
	if err := mgr.storage.Find(&resources, "id IN (?)", ids); err != nil {
		return nil, errors.Trace(err)
	}
	return resources, nil
}

// GetResourceByHost get resource by host
func (mgr *Manager) GetResourceByHost(host string) (*types.Resource, error) {
	var resource types.Resource
	if err := mgr.storage.FindOne(&resource, "host=?", host); err != nil {
		return nil, errors.NotFoundf("resource host %d", host)
	}
	return &resource, nil
}

// UpdateResourceOne update one resource
func (mgr *Manager) UpdateResourceOne(resource *types.Resource) error {
	return errors.Trace(mgr.storage.Save(resource))
}

// FreeResource free resources in range
func (mgr *Manager) LockResource(cond []int) error {
	return errors.Trace(mgr.storage.Update(&types.Resource{}, "id IN (?)", []interface{}{cond},
		map[string]interface{}{"status": types.ResourceStatusWait}))
}

// FreeResource free resources in range
func (mgr *Manager) FreeResource(cond []int) error {
	return errors.Trace(mgr.storage.Update(&types.Resource{}, "id IN (?)", []interface{}{cond},
		map[string]interface{}{"status": types.ResourceStatusReady}))
}

// ClearResourceLock free all resources
// IT MAY CAUSE RESOURCE COLLISION
// DO NOT USE IT ANYWHERE EXCEPT SCHEDULER START
func (mgr *Manager) ClearResourceLock() error {
	return errors.Trace(mgr.storage.UpdateAll(&types.Resource{},
		map[string]interface{}{"status": types.ResourceStatusReady}))
}
