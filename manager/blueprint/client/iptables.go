package client

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/you06/chaos-commander/manager/executor/content"
	"github.com/you06/chaos-commander/pkg/types"
)

type IpTablesClient struct {
	sshClient *sshClient
	resource  *types.Resource
}

func newIpTables(resource *types.Resource) *IpTablesClient {
	client := IpTablesClient{
		resource: resource,
	}
	return &client
}

func (i *IpTablesClient) drop(addrStr string) error {
	client, err := newSSHClient(i.resource)
	if err != nil {
		return errors.Trace(err)
	}
	i.sshClient = client

	cmd := fmt.Sprintf("iptables -A INPUT -s %s -j DROP", addrStr)
	err = i.sshClient.Run(cmd)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (i *IpTablesClient) dropClear(addrStr string) error {
	client, err := newSSHClient(i.resource)
	if err != nil {
		return errors.Trace(err)
	}
	i.sshClient = client

	cmd := fmt.Sprintf("for i in `iptables -S | grep \"\\-A INPUT \\-s %s.*\\-j DROP\" | wc -l | xargs seq 1`; do iptables -D INPUT -s %s -j DROP; done",
		addrStr, addrStr)
	err = i.sshClient.Run(cmd)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (i *IpTablesClient) reject(addrStr string) error {
	client, err := newSSHClient(i.resource)
	if err != nil {
		return errors.Trace(err)
	}
	i.sshClient = client

	cmd := fmt.Sprintf("iptables -A INPUT -s %s -p TCP -j REJECT --reject-with tcp-reset", addrStr)
	err = i.sshClient.Run(cmd)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (i *IpTablesClient) rejectClear(addrStr string) error {
	client, err := newSSHClient(i.resource)
	if err != nil {
		return errors.Trace(err)
	}
	i.sshClient = client

	cmd := fmt.Sprintf("for i in `iptables -S | grep \"\\-A INPUT \\-s %s.*\\-j REJECT\" | wc -l | xargs seq 1`; do iptables -D INPUT -s %s -p TCP -j REJECT --reject-with tcp-reset; done",
		addrStr, addrStr)
	err = i.sshClient.Run(cmd)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func IpTablesDrop(ctn *content.Content, resource *types.Resource, addrStr string, args ...interface{}) error {
	i := newIpTables(resource)
	return errors.Trace(i.drop(addrStr))
}

func IpTablesDropClear(ctn *content.Content, resource *types.Resource, addrStr string, args ...interface{}) error {
	i := newIpTables(resource)
	return errors.Trace(i.dropClear(addrStr))
}

func IpTablesReject(ctn *content.Content, resource *types.Resource, addrStr string, args ...interface{}) error {
	i := newIpTables(resource)
	return errors.Trace(i.reject(addrStr))
}

func IpTablesRejectClear(ctn *content.Content, resource *types.Resource, addrStr string, args ...interface{}) error {
	i := newIpTables(resource)
	return errors.Trace(i.rejectClear(addrStr))
}
