package content

import (
	"github.com/you06/chaos-commander/pkg"
	"github.com/you06/chaos-commander/pkg/log"
	"github.com/you06/chaos-commander/pkg/types"
)

type Content struct{
	Pkg       *pkg.Pkg
	Job       *types.Job
	Ansible   *types.Ansible
	Blueprint *types.BluePrint
	Load      *types.Load
	Resources []*types.Resource
	Logger    *log.Log
}
