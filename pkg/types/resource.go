package types

type CGroupType int
type ResourceStatus int

const (
	_ CGroupType = iota
	CGroupV1
	CGroupV2
)

const (
	ResourceStatusReady ResourceStatus = iota
	ResourceStatusWait
)

type Resource struct {
	ID       int            `gorm:"column:id"`
	Host     string         `gorm:"column:host"`
	Port     int            `gorm:"column:port"`
	CGroup   CGroupType     `gorm:"column:cgroup"`
	Password string         `gorm:"column:password"`
	Key      string         `gorm:"column:key"`
	Status   ResourceStatus `gorm:"column:status"`
}

func (r *Resource) GetID() int {
	return r.ID
}

func (r *Resource) GetHost() string {
	return r.Host
}

func (r *Resource) GetPort() int {
	return r.Port
}

func (r *Resource) GetPassword() string {
	return r.Password
}

func (r *Resource) GetKey() string {
	return r.Key
}

func (r *Resource) GetCGroup() CGroupType {
	return r.CGroup
}

func (r *Resource) GetStatus() ResourceStatus {
	return r.Status
}

func (c CGroupType) String() string {
	switch c {
	case CGroupV1:
		return "CGroup V1"
	case CGroupV2:
		return "CGroup V2"
	default:
		return "unknown CGroup version"
	}
}

func (s ResourceStatus) String() string {
	switch s {
	case ResourceStatusReady:
		return "ready"
	case ResourceStatusWait:
		return "wait"
	default:
		return "unknown resource status"
	}
}
