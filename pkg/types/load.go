package types

type LoadType int

const (
	_ LoadType = iota
	LoadTypePath
	LoadTypeGit
)

type Load struct {
	ID    int      `gorm:"column:id"`
	Type  LoadType `gorm:"column:type"`
	Path  string   `gorm:"column:path"`
	Git   string   `gorm:"column:git"`
	Param string   `gorm:"column:param"`
}

func (l *Load) GetID() int {
	return l.ID
}

func (l *Load) GetType() LoadType {
	return l.Type
}

func (l *Load) GetPath() string {
	return l.Path
}

func (l *Load) GetGit() string {
	return l.Git
}

func (l *Load) GetParam() string {
	return l.Param
}

func (t LoadType) String() string {
	switch t {
	case LoadTypePath:
		return "load type path"
	case LoadTypeGit:
		return "load type git"
	default:
		return "unknown load type"
	}
}
