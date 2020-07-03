package client

import (
	"fmt"
	"os/user"
	"path"
	"strings"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/manager/executor/content"
	"github.com/you06/chaos-commander/pkg/types"
	"github.com/you06/chaos-commander/util"
)

type SuddenDeath struct {
	resource *types.Resource
}

func newSuddenDeath(resource *types.Resource) *SuddenDeath {
	client := SuddenDeath{
		resource: resource,
	}
	return &client
}

func (s *SuddenDeath) suddenDeathShutdown() error {
	usr, err := user.Current()
	if err != nil {
		return errors.Trace(err)
	}
	controlResource := &types.Resource{
		Host:   "127.0.0.1", // address of central control machine
		Port:   22,
		CGroup: types.CGroupV1,
		Key:    path.Join(usr.HomeDir, ".ssh/id_rsa"),
	}
	client, err := newSSHClient(controlResource)
	if err != nil {
		return errors.Trace(err)
	}

	host := strings.Replace(s.resource.GetHost(), "172.16", "192.168", 1)

	suddenDeathCmd := fmt.Sprintf("/opt/dell/srvadmin/bin/idracadm7 -r %s -u root -p 'password' serveraction powercycle", host)

	log.Info(suddenDeathCmd)
	err = client.Run(suddenDeathCmd)
	if err != nil {
		log.Error(errors.ErrorStack(err))
		return errors.Trace(err)
	}

	return nil
}

func SuddenDeathDo(ctn *content.Content, resource *types.Resource, args ...interface{}) error {
	s := newSuddenDeath(resource)
	return errors.Trace(s.suddenDeathShutdown())
}

func SuddenDeathClear(ctn *content.Content, resource *types.Resource, args ...interface{}) error {
	// start service by calling ansible
	_, err := util.DoCmd(ctn.Ansible.GetPath(), "ansible-playbook", []string{
		"start.yml",
		"-l",
		resource.GetHost(),
	}...)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
