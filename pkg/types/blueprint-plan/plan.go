package blueprint_plan

// TODO: separate them into individual files

// BluePrintPlan defines blue print plan
type BluePrintPlan struct {
	Hosts      []*Host      `json:"hosts" yaml:"hosts"`
	Regions    []*Region    `json:"regions" yaml:"regions"`
	Relations  []*Relation  `json:"relations" yaml:"relations"`
	Exceptions []*Exception `json:"exceptions" yaml:"exceptions"`
	Steps      []*Step      `json:"steps" yaml:"steps"`
}

// Host defines host structure
type Host struct {
	Name  string `json:"name" yaml:"name"`
	Host  string `json:"host" yaml:"host"`
	Ether string `json:"ether" yaml:"ether"`
}

// Region defines region structure
type Region struct {
	Name  string   `json:"name" yaml:"name"`
	Nodes []string `json:"nodes" yaml:"nodes"`
}

// Relation defines relation structure
type Relation struct {
	Name string `json:"name" yaml:"name"`
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
	Rule string `json:"rule" yaml:"rule"`
}

// Exception defines exception structure
type Exception struct {
	Name string `json:"name" yaml:"name"`
	Host string `json:"host" yaml:"host"`
	Rule string `json:"rule" yaml:"rule"`
}

// Step defines step structure
type Step struct {
	Name     string        `json:"name" yaml:"name"`
	Duration int           `json:"duration" yaml:"duration"`
	Args     []interface{} `json:"args" yaml:"args"`
}

func (h *Host) GetName() string {
	return h.Name
}

func (h *Host) GetHost() string {
	return h.Host
}

func (h *Host) GetEther() string {
	return h.Ether
}

func (r *Region) GetName() string {
	return r.Name
}

func (s *Step) GetName() string {
	return s.Name
}

func (s *Step) GetDuration() int {
	return s.Duration
}

func (s *Step) GetArgs() []interface{} {
	return s.Args
}

func (e *Exception) GetName() string {
	return e.Name
}

func (e *Exception) GetHost() string {
	return e.Host
}

func (e *Exception) GetRule() string {
	return e.Rule
}

func (r *Relation) GetName() string {
	return r.Name
}

func (r *Relation) GetRule() string {
	return r.Rule
}
