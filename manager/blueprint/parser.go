package blueprint

import (
	"encoding/json"
	"github.com/juju/errors"
	plan "github.com/you06/chaos-commander/pkg/types/blueprint-plan"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func (b *Blueprint)parse() error {
	var (
		p plan.BluePrintPlan
	)

	if b.content.Blueprint.GetPath() != "" {
		blueprintContent, err := ioutil.ReadFile(b.content.Blueprint.GetPath())
		if err != nil {
			return errors.Trace(err)
		}
		if err := yaml.Unmarshal(blueprintContent, &p); err != nil {
			return errors.Trace(err)
		}
	} else {
		if err := json.Unmarshal([]byte(b.content.Blueprint.GetBluePrint()), &p); err != nil {
			return errors.Trace(err)
		}
	}
	b.plan = &p

	if err := b.parseHost(); err != nil {
		return errors.Trace(err)
	}
	if err := b.checkResource(); err != nil {
		return errors.Trace(err)
	}
	if err := b.parseRegion(); err != nil {
		return errors.Trace(err)
	}
	if err := b.parseRelation(); err != nil {
		return errors.Trace(err)
	}
	if err := b.parseException(); err != nil {
		return errors.Trace(err)
	}
	if err := b.parseStep(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (b *Blueprint)parseHost() error {
	b.hosts = make(map[string]*plan.Host)
	for _, host := range b.plan.Hosts {
		if _, exist := b.hosts[host.GetName()]; exist {
			return errors.New("duplicated host")
		}
		b.hosts[host.GetName()] = host
	}
	return nil
}

func (b *Blueprint)checkResource() error {
	for hostname, _ := range b.hosts {
		host := b.hosts[hostname].GetHost()
		if _, exist := b.resourceMap[host]; !exist {
			return errors.NotFoundf("host %s not found", host)
		}
	}
	return nil
}

func (b *Blueprint)parseRegion() error {
	b.regions = make(map[string]*plan.Region)
	for _, region := range b.plan.Regions {
		if _, exist := b.hosts[region.GetName()]; exist {
			return errors.New("region name duplicated with host name")
		}
		if _, exist := b.regions[region.GetName()]; exist {
			return errors.New("duplicated host")
		}
		b.regions[region.GetName()] = region
	}
	return nil
}

// TODO: check if the input valid
func (b *Blueprint)parseRelation() error {
	var relations []*plan.Relation
	for _, relation := range b.plan.Relations {
		relations = append(relations, relation)
	}
	b.relations = relations
	return nil
}

// TODO: check if the input valid
func (b *Blueprint)parseException() error {
	var exceptions []*plan.Exception
	for _, exception := range b.plan.Exceptions {
		exceptions = append(exceptions, exception)
	}
	b.exceptions = exceptions
	return nil
}

// TODO: check if the input valid
func (b *Blueprint)parseStep() error {
	var steps []*plan.Step
	for _, step := range b.plan.Steps {
		steps = append(steps, step)
	}
	b.steps = steps
	return nil
}

func (b *Blueprint)getHost(name string) (*plan.Host, error) {
	host, exist := b.hosts[name]
	if exist {
		return host, nil
	}
	return nil, errors.NotFoundf("host %s", name)
}

func (b *Blueprint)getRegion(name string) ([]*plan.Host, error) {
	var hosts []*plan.Host
	region, exist := b.regions[name]
	if exist {
		for _, node := range region.Nodes {
			host, err := b.getHost(node)
			if err != nil {
				return nil, errors.Trace(err)
			}
			hosts = append(hosts, host)
		}
	}
	host, exist := b.hosts[name]
	if exist {
		hosts = append(hosts, host)
	}

	if len(hosts) > 0 {
		return hosts, nil
	}
	return nil, errors.NotFoundf("region %s", name)
}

func (b *Blueprint)getRelations(name string) []*plan.Relation {
	var relations []*plan.Relation
	for _, relation := range b.relations {
		if relation.GetName() == name {
			relations = append(relations, relation)
		}
	}
	return relations
}

func (b *Blueprint)getExceptions(name string) []*plan.Exception {
	var exceptions []*plan.Exception
	for _, exception := range b.exceptions {
		if exception.GetName() == name {
			exceptions = append(exceptions, exception)
		}
	}
	return exceptions
}
