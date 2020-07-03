package blueprint

import (
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/manager/executor/content"
	"github.com/you06/chaos-commander/pkg/types"
	plan "github.com/you06/chaos-commander/pkg/types/blueprint-plan"
)

// Blueprint defines blueprint
type Blueprint struct {
	content      *content.Content
	relationMap  relationMap
	exceptionMap exceptionMap

	plan       *plan.BluePrintPlan
	hosts      map[string]*plan.Host
	regions    map[string]*plan.Region
	relations  []*plan.Relation
	exceptions []*plan.Exception
	steps      []*plan.Step
	stepch     chan *plan.Step

	resourceMap   map[string]*types.Resource

	finish bool
	callback func()
}

// New creates a blueprint
func New(content *content.Content, resources []*types.Resource, callback func()) (*Blueprint, error) {
	resourceMap := make(map[string]*types.Resource)
	for _, resource := range resources {
		k := resource.GetHost()
		resourceMap[k] = resource
	}

	b := Blueprint{
		content: content,
		relationMap: initRelationMap(),
		exceptionMap: initExceptionMap(),
		stepch: make(chan *plan.Step, 1),
		resourceMap: resourceMap,
		finish: false,
		callback: callback,
	}

	if err := b.parse(); err != nil {
		return nil, errors.Trace(err)
	}

	return &b, nil
}

// Start runs the chaos plan
func (b *Blueprint)Start() {
	log.Infof("job: %d, blueprint: %d start", b.content.Job.GetID(), b.content.Blueprint.GetID())
	b.content.Logger.Infof("job: %d, blueprint: %d start", b.content.Job.GetID(), b.content.Blueprint.GetID())
	go b.exec()
	for _, step := range b.steps {
		b.stepch <- step
	}
	close(b.stepch)
}
