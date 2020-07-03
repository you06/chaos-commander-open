package types

type Task struct {
	Job       *Job
	Ansible   *Ansible
	BluePrint *BluePrint
	Load      *Load
	Resources []*Resource
}

func (t *Task) GetJob() *Job {
	if t.Job != nil {
		return t.Job
	}
	return &Job{}
}

func (t *Task) GetAnsible() *Ansible {
	if t.Ansible != nil {
		return t.Ansible
	}
	return &Ansible{}
}

func (t *Task) GetBluePrint() *BluePrint {
	if t.BluePrint != nil {
		return t.BluePrint
	}
	return &BluePrint{}
}

func (t *Task) GetLoad() *Load {
	if t.Load != nil {
		return t.Load
	}
	return &Load{}
}

func (t *Task) GetResources() []*Resource {
	return t.Resources
}
