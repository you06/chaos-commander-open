package types

type Ansible struct {
	ID     int    `gorm:"column:id"`
	User   string `gorm:"column:user"`
	Path   string `gorm:"column:path"`
	Deploy string `gorm:"column:deploy"`
	Prom   string `gorm:"column:prom"`
}

func (a *Ansible) GetID() int {
	return a.ID
}

func (a *Ansible) GetUser() string {
	return a.User
}

func (a *Ansible) GetPath() string {
	return a.Path
}

func (a *Ansible) GetDeploy() string {
	return a.Deploy
}

func (a *Ansible) GetProm() string {
	return a.Prom
}
