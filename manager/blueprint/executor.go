package blueprint

import (
	"github.com/ngaut/log"
	plan "github.com/you06/chaos-commander/pkg/types/blueprint-plan"
	"time"
)

const MIN_STEP_SLEEP = 5

type relationFn func(fromHost, toHost *plan.Host)
type exceptionFn func(host *plan.Host)

func (b *Blueprint)exec() {
	for {
		step := <- b.stepch
		if b.finish {
			continue
		}
		if step != nil {
			b.doStep(step)
		} else {
			log.Infof("job: %d, blueprint: %d finish", b.content.Job.GetID(), b.content.Blueprint.GetID())
			b.content.Logger.Infof("job: %d, blueprint: %d finish", b.content.Job.GetID(), b.content.Blueprint.GetID())
			b.finish = true
			b.callback()
			break
		}
	}
}

func (b *Blueprint)doStep(step *plan.Step) {
	if step.GetName() == "sleep"  {
		b.content.Logger.Infof("sleep for %d seconds", step.GetDuration())
		time.Sleep(time.Duration(step.GetDuration()) * time.Second)
		return
	}
	relations := b.getRelations(step.GetName())
	exceptions := b.getExceptions(step.GetName())
	go b.doRelations(relations, step)
	go b.doExceptions(exceptions, step)
	duration := step.GetDuration()
	if duration < MIN_STEP_SLEEP {
		duration = MIN_STEP_SLEEP
	}
	time.Sleep(time.Duration(duration) * time.Second)
	b.clearRelations(relations, step)
	b.clearExceptions(exceptions, step)
	time.Sleep(time.Duration(MIN_STEP_SLEEP) * time.Second)
}

func (b *Blueprint)mapRelation(fromRegionName, toRegionName string, fn relationFn) {
	fromRegion, err := b.getRegion(fromRegionName)
	if err != nil {
		log.Errorf("get from region err %v", err)
	}
	toRegion, err := b.getRegion(toRegionName)
	if err != nil {
		log.Errorf("get from region err %v", err)
	}
	for _, fromHost := range fromRegion {
		for _, toHost := range toRegion {
			fn(fromHost, toHost)
		}
	}
}

func (b *Blueprint)mapException(regionName string, fn exceptionFn) {
	region, err := b.getRegion(regionName)
	if err != nil {
		log.Errorf("get region err %v", err)
	}
	for _, host := range region {
		fn(host)
	}
}

func (b *Blueprint)doRelations(relations []*plan.Relation, step *plan.Step) {
	for _, relation := range relations {
		b.mapRelation(relation.From, relation.To, func(fromHost, toHost *plan.Host) {
			err := b.doRelation(fromHost, toHost, relation, step)
			if err != nil {
				log.Errorf("do relation err %v, %v", err, relation)
			}
		})
	}
}

func (b *Blueprint)doExceptions(exceptions []*plan.Exception, step *plan.Step) {
	for _, exception := range exceptions {
		b.mapException(exception.GetHost(), func(host *plan.Host) {
			err := b.doException(host, exception, step)
			if err != nil {
				log.Errorf("do exception err %v, %v", err, exception)
			}
		})
	}
}

func (b *Blueprint)clearRelations(relations []*plan.Relation, step *plan.Step) {
	for _, relation := range relations {
		b.mapRelation(relation.From, relation.To, func(fromHost, toHost *plan.Host) {
			err := b.clearRelation(fromHost, toHost, relation, step)
			if err != nil {
				log.Errorf("clear relation err %v, %v", err, relation)
			}
		})
	}
}

func (b *Blueprint)clearExceptions(exceptions []*plan.Exception, step *plan.Step) {
	for _, exception := range exceptions {
		b.mapException(exception.GetHost(), func(host *plan.Host) {
			err := b.clearException(host, exception, step)
			if err != nil {
				log.Errorf("clear exception err %v, %v", err, exception)
			}
		})
	}
}
