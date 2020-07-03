package types

import "github.com/you06/chaos-commander/util"

type Job struct {
	ID          int    `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	RunningPlan string `gorm:"column:running_plan"`
	Tidb        string `gorm:"column:tidb"`
	Tikv        string `gorm:"column:tikv"`
	Pd          string `gorm:"column:pd"`
	Resources   string `gorm:"column:resources"`
	AnsibleID   int    `gorm:"column:ansible"`
	BluePrintID int    `gorm:"column:blue_print"`
	LoadID      int    `gorm:"column:load"`
	SkipDeploy  bool   `gorm:"column:skip_deploy"`
	Pessimistic bool   `gorm:"column:pessimistic"`
	DSN         string `gorm:"column:dsn"`
}

func (j *Job) GetID() int {
	return j.ID
}

func (j *Job) GetName() string {
	return j.Name
}

func (j *Job) GetRunningPlan() string {
	return j.RunningPlan
}

func (j *Job) GetTidb() string {
	return j.Tidb
}

func (j *Job) GetTikv() string {
	return j.Tikv
}

func (j *Job) GetPd() string {
	return j.Pd
}

func (j *Job) GetAnsibleID() int {
	return j.AnsibleID
}

func (j *Job) GetBluePrintID() int {
	return j.BluePrintID
}

func (j *Job) GetLoadID() int {
	return j.LoadID
}

func (j *Job) GetResources() []int {
	return util.ParseIntSlice(j.Resources)
}

func (j *Job) GetSkipDeploy() bool {
	return j.SkipDeploy
}

func (j *Job) GetPessimistic() bool {
	return j.Pessimistic
}

func (j *Job) GetDSN() string {
	return j.DSN
}
