package client

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/you06/chaos-commander/manager/executor/content"
	"github.com/you06/chaos-commander/pkg/types"
)

type CGroupV2Client struct {
	sshClient *sshClient
	resource  *types.Resource
}

func newCGroupV2(resource *types.Resource) *CGroupV2Client {
	client := CGroupV2Client{
		resource: resource,
	}
	return &client
}

func (c *CGroupV2Client) cGroupV2ioMax(dev, name string, speed int64, limitType string) error {
	client, err := newSSHClient(c.resource)
	if err != nil {
		return errors.Trace(err)
	}
	c.sshClient = client

	processID, err := c.getProcessID(name)
	if err != nil {
		return errors.Trace(err)
	}
	diskID, err := c.getDiskID(dev)
	if err != nil {
		return errors.Trace(err)
	}

	ioGroupPath := fmt.Sprintf("/sys/fs/cgroup/unified/%s/io.max", GROUP_NAME)
	ioCommand := fmt.Sprintf("echo \"%s %s=%d\" > %s", diskID, limitType, speed, ioGroupPath)
	pidCommand := fmt.Sprintf("echo %s > /sys/fs/cgroup/unified/%s/cgroup.procs", processID, GROUP_NAME)

	err = c.sshClient.Run(ioCommand)
	if err != nil {
		return errors.Trace(err)
	}
	err = c.sshClient.Run(pidCommand)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *CGroupV2Client) cGroupV2ioClear(dev string) error {
	client, err := newSSHClient(c.resource)
	if err != nil {
		return errors.Trace(err)
	}
	c.sshClient = client

	diskID, err := c.getDiskID(dev)
	if err != nil {
		return errors.Trace(err)
	}
	ioGroupPath := fmt.Sprintf("/sys/fs/cgroup/unified/%s/io.max", GROUP_NAME)
	clearCmd := fmt.Sprintf("echo \"%s rbps=max wbps=max riops=max wiops=max\" > %s", diskID, ioGroupPath)

	err = c.sshClient.Run(clearCmd)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *CGroupV2Client) getProcessID(name string) (string, error) {
	cmd := fmt.Sprintf("ps -ef| grep '%s' | grep -v grep | awk '{print $2}'", name)
	out, err := c.sshClient.RunOutput(cmd)
	if err != nil {
		return "", errors.Trace(err)
	}
	return out, nil
}

func (c *CGroupV2Client) getDiskID(dev string) (string, error) {
	cmd := fmt.Sprintf("lsblk | grep '%s' | head -n 1 | awk '{print $2}'", dev)
	out, err := c.sshClient.RunOutput(cmd)
	if err != nil {
		return "", errors.Trace(err)
	}
	return out, nil
}

func CGroupV2ioWbps(ctn *content.Content, resource *types.Resource, args ...interface{}) error {
	devStr, ok := args[0].(string)
	if !ok {
		return errors.New("dev not string")
	}
	nameStr, ok := args[1].(string)
	if !ok {
		return errors.New("name not string")
	}

	var speedInt64 int64
	speedFloat, okFloat := args[2].(float64)
	speedInt, okInt := args[2].(int)
	if okFloat {
		speedInt64 = int64(speedFloat)
	} else if okInt {
		speedInt64 = int64(speedInt)
	} else {
		return errors.New("speed not a number")
	}

	c := newCGroupV2(resource)
	return errors.Trace(c.cGroupV2ioMax(devStr, nameStr, speedInt64, "wbps"))
}

func CGroupV2ioClear(ctn *content.Content, resource *types.Resource, args ...interface{}) error {
	devStr, ok := args[0].(string)
	if !ok {
		return errors.New("dev not string")
	}

	c := newCGroupV2(resource)
	return errors.Trace(c.cGroupV2ioClear(devStr))
}
