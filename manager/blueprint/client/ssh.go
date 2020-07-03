package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/types"
	"golang.org/x/crypto/ssh"
)

type sshClient struct {
	client *ssh.Client
}

func newSSHClient(resource *types.Resource) (*sshClient, error) {
	var auths []ssh.AuthMethod
	if resource.GetKey() != "" {
		if auth, err := publicKeyFile(resource.GetKey()); err != nil {
			return nil, errors.Trace(err)
		} else {
			auths = append(auths, auth)
		}
	} else if resource.GetPassword() != "" {
		auths = append(auths, ssh.Password(resource.GetPassword()))
	}

	config := ssh.ClientConfig{
		User:            "root",
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp",
		fmt.Sprintf("%s:%d", resource.GetHost(), resource.GetPort()),
		&config)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return &sshClient{
		client,
	}, nil
}

func (s *sshClient) RunOutput(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", errors.Trace(err)
	}
	defer session.Close()
	out, err := session.Output(cmd)
	if err != nil {
		return "", errors.Trace(err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *sshClient) Run(cmd string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return errors.Trace(err)
	}
	defer session.Close()
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Println(string(out))
		return errors.Trace(err)
	}
	return nil
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Trace(err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return ssh.PublicKeys(key), nil
}

func NewSSHClient(resource *types.Resource) (*sshClient, error) {
	return newSSHClient(resource)
}
