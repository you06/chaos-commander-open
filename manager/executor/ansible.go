package executor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/you06/chaos-commander/pkg/types"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/manager/blueprint/client"
	"github.com/you06/chaos-commander/util"

	"path"
	"strings"
	"sync"
)

const (
	FILE_SERVER_URL = "http://fileserver.pingcap.net"
	SHA_TEMPLATE    = "%s/download/refs/pingcap/%s/%s/sha1"
	URL_TEMPLATE    = "%s/download/builds/pingcap/%s/%s/centos7/tidb-server.tar.gz"
)

func (e *Executor) ansibleDownload() error {
	if err := e.ansibleClear(); err != nil {
		return errors.Trace(err)
	}
	if err := e.ansibleDownloadComponents(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *Executor) ansibleClear() error {
	out, err := util.DoCmd(e.content.Ansible.GetPath(), "rm", "-rf", "downloads")
	if err != nil {
		e.logger.Info(out)
		return errors.Trace(err)
	}
	e.logger.Info("clear downloads directory")
	out, err = util.DoCmd(e.content.Ansible.GetPath(), "rm", "-rf", "resources")
	if err != nil {
		e.logger.Info(out)
		return errors.Trace(err)
	}
	e.logger.Info("clear resources directory")
	return nil
}

func (e *Executor) ansibleDownloadComponents() error {
	out := ""
	err := util.RetryOnError(context.Background(), 3, func() error {
		o, err := util.DoCmd(e.content.Ansible.GetPath(), "ansible-playbook", "local_prepare.yml")
		out = o
		return errors.Trace(err)
	})
	if err != nil {
		e.logger.Info("ansible download components failed")
		e.logger.Info(out)
		return errors.Trace(err)
	} else {
		e.logger.Info("ansible download components success")
	}
	return nil
}

func (e *Executor) download() {
	var wg sync.WaitGroup
	wg.Add(3)
	go e.downloadComponent(&wg, "tidb")
	go e.downloadComponent(&wg, "tikv")
	go e.downloadComponent(&wg, "pd")
	wg.Wait()
}

func (e *Executor) downloadComponent(wg *sync.WaitGroup, component string) {
	defer wg.Done()
	switch component {
	case "tidb":
		e.downloadTiDB()
	case "tikv":
		e.downloadTiKV()
	case "pd":
		e.downloadPD()
	}
}

func (e *Executor) downloadTiDB() {
	sha, err := e.getSHA(fmt.Sprintf(SHA_TEMPLATE, FILE_SERVER_URL, "tidb", e.content.Job.GetTidb()))
	if err != nil {
		log.Errorf("tidb get sha err %v", err)
	}
	log.Info(sha)
}

func (e *Executor) downloadTiKV() {
	sha, err := e.getSHA(fmt.Sprintf(SHA_TEMPLATE, FILE_SERVER_URL, "tikv", e.content.Job.GetTidb()))
	if err != nil {
		log.Errorf("tidb get sha err %v", err)
	}
	log.Info(sha)
}

func (e *Executor) downloadPD() {
	sha, err := e.getSHA(fmt.Sprintf(SHA_TEMPLATE, FILE_SERVER_URL, "pd", e.content.Job.GetTidb()))
	if err != nil {
		log.Errorf("tidb get sha err %v", err)
	}
	log.Info(sha)
}

func (e *Executor) getSHA(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Trace(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	return strings.TrimSpace(string(body)), nil
}

func (e *Executor) FetchBinaryAndExtract(url string) {

}

func (e *Executor) ansibleStop() error {
	out := ""
	err := util.RetryOnError(context.Background(), 3, func() error {
		o, err := util.DoCmd(e.content.Ansible.GetPath(), "ansible-playbook", "stop.yml")
		out = o
		return errors.Trace(err)
	})
	if err != nil {
		e.logger.Info("ansible stop failed")
		e.logger.Info(out)
		return errors.Trace(err)
	} else {
		e.logger.Info("ansible stop success")
	}
	return nil
}

func (e *Executor) ansibleClearLog() error {
	// no log clear workload in ansible
	// using ssh to clear log :(
	nodes := len(e.content.Resources)
	var wg sync.WaitGroup
	wg.Add(len(e.content.Resources))
	remoteLogPath := fmt.Sprintf("%s/*", path.Join(e.content.Ansible.GetDeploy(), "log"))
	for _, resource := range e.content.Resources {
		go func(resource *types.Resource) {
			defer wg.Done()
			sshClient, err := client.NewSSHClient(resource)
			if err != nil {
				e.logger.Errorf("ssh %s error, %v", resource.GetHost(), err)
				return
			}
			e.logger.Infof("deploy log path is %s", remoteLogPath)
			if err := sshClient.Run(fmt.Sprintf("rm -rf %s", remoteLogPath)); err != nil {
				e.logger.Errorf("remove expired logs on %s failed, %v", resource.GetHost())
			} else {
				e.logger.Infof("remove expired logs on %s success", resource.GetHost())
				nodes--
			}
		}(resource)
	}
	wg.Wait()
	if nodes == 0 {
		return nil
	}
	return errors.Errorf("%d nodes log clear failed", nodes)
}

func (e *Executor) ansibleNTP() error {
	out := ""
	err := util.RetryOnError(context.Background(), 3, func() error {
		o, err := util.DoCmd(e.content.Ansible.GetPath(), "ansible-playbook", []string{
			"-i",
			"hosts.ini",
			"deploy_ntp.yml",
			"-u",
			e.content.Ansible.GetUser(),
			"-b",
		}...)
		out = o
		return errors.Trace(err)
	})
	if err != nil {
		e.logger.Info("ansible deploy ntp failed")
		e.logger.Info(out)
		return errors.Trace(err)
	}

	e.logger.Info("ansible deploy ntp success")
	return nil
}

func (e *Executor) ansibleCPU() error {
	err := util.RetryOnError(context.Background(), 3, func() error {
		_, err := util.DoCmd(e.content.Ansible.GetPath(), "ansible", []string{
			"-i",
			"hosts.ini",
			"all",
			"-m",
			"shell",
			"-a",
			"cpupower frequency-set --governor performance",
			"-u",
			e.content.Ansible.GetUser(),
			"-b",
		}...)
		return errors.Trace(err)
	})
	// error may caused by VM which don't support CPU governor
	if err != nil {
		e.logger.Info("cpu performance failed")
		//return errors.Trace(err)
		return nil
	}

	e.logger.Info("cpu performance success")
	return nil
}

func (e *Executor) ansibleSwapOff() error {
	out := ""
	err := util.RetryOnError(context.Background(), 3, func() error {
		o, err := util.DoCmd(e.content.Ansible.GetPath(), "ansible", []string{
			"-i",
			"hosts.ini",
			"all",
			"-m",
			"shell",
			"-a",
			"swapoff -a",
			"-u",
			e.content.Ansible.GetUser(),
			"-b",
		}...)
		out = o
		return errors.Trace(err)
	})

	if err != nil {
		e.logger.Error("ansible swapoff failed")
		e.logger.Error(out)
		return errors.Trace(err)
	}

	e.logger.Info("ansible swapoff success")
	return nil
}

func (e *Executor) ansibleDeploy() error {
	out := ""
	err := util.RetryOnError(context.Background(), 3, func() error {
		o, err := util.DoCmd(e.content.Ansible.GetPath(), "ansible-playbook", "deploy.yml")
		out = o
		return errors.Trace(err)
	})
	if err != nil {
		e.logger.Info("ansible deploy failed")
		e.logger.Info(out)
		return errors.Trace(err)
	}

	e.logger.Info("ansible deploy success")
	return nil
}

func (e *Executor) ansibleStart() error {
	out := ""
	err := util.RetryOnError(context.Background(), 3, func() error {
		//	out, err := util.DoCmd(e.content.Ansible.GetPath(), "ls")
		o, err := util.DoCmd(e.content.Ansible.GetPath(), "ansible-playbook", "start.yml")
		out = o
		return errors.Trace(err)
	})
	if err != nil {
		e.logger.Info("ansible start failed")
		e.logger.Info(out)
		return errors.Trace(err)
	}

	e.logger.Info("ansible start success")
	return nil
}
