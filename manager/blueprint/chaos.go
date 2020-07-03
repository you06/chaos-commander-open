package blueprint

import (
	"github.com/juju/errors"
	plan "github.com/you06/chaos-commander/pkg/types/blueprint-plan"
)


func (b *Blueprint)doRelation(fromHost, toHost *plan.Host, relation *plan.Relation, step *plan.Step) error {
	r, exist := b.relationMap[relation.GetRule()]
	if !exist {
		return errors.NotFoundf("%s rule not found", relation.GetRule())
	}
	b.content.Logger.Infof("do relation %s on from %s to %s", relation.GetRule(), fromHost.GetHost(), toHost.GetHost())
	resource := b.resourceMap[fromHost.GetHost()]
	err := r.Do.Call(b.content, resource, toHost.GetHost(), step.Args)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (b *Blueprint)doException(host *plan.Host, exception *plan.Exception, step *plan.Step) error {
	e, exist := b.exceptionMap[exception.GetRule()]
	if !exist {
		return errors.NotFoundf("%s rule not found", exception.GetRule())
	}
	b.content.Logger.Infof("do exception %s on %s", exception.GetRule(), host.GetHost())
	resource := b.resourceMap[host.GetHost()]
	err := e.Do.Call(b.content, resource, step.Args)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (b *Blueprint)clearRelation(fromHost, toHost *plan.Host, relation *plan.Relation, step *plan.Step) error {
	r, exist := b.relationMap[relation.GetRule()]
	if !exist {
		return errors.NotFoundf("%s rule not found", relation.GetRule())
	}
	b.content.Logger.Infof("clear relation %s from %s to %s", relation.GetRule(), fromHost.GetHost(), toHost.GetHost())
	resource := b.resourceMap[fromHost.GetHost()]
	err := r.Clear.Call(b.content, resource, toHost.GetHost(), step.Args)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *Blueprint)clearException(host *plan.Host, exception *plan.Exception, step *plan.Step) error {
	e, exist := b.exceptionMap[exception.GetRule()]
	if !exist {
		return errors.NotFoundf("%s rule not found", exception.GetRule())
	}
	b.content.Logger.Infof("clear exception %s on %s", exception.GetRule(), host.GetHost())
	resource := b.resourceMap[host.GetHost()]
	err := e.Clear.Call(b.content, resource, step.Args)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

