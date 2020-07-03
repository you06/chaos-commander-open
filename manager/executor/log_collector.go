package executor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/juju/errors"
	nlog "github.com/ngaut/log"
	"github.com/you06/chaos-commander/pkg/log"
	"github.com/you06/chaos-commander/pkg/types"
	"github.com/you06/chaos-commander/util"
)

type logAnalyse struct {
	info  int
	warn  int
	error int
	fatal int
}

func (e *Executor) fetchAllLog() {
	e.content.Logger.Info("fetch all log start")
	message := ""
	var wg sync.WaitGroup
	for _, resource := range e.content.Resources {
		wg.Add(1)
		go e.fetchLog(&wg, resource, &message)
	}
	wg.Wait()
	title := fmt.Sprintf("Job %d, History %d", e.history.GetJobID(), e.history.GetID())

	if err := e.mgr.Pkg.Slack.SendMessage(fmt.Sprintf("%s\n```%s```", title, message)); err != nil {
		nlog.Errorf("slack msg error %v", errors.ErrorStack(err))
	}
}

// TODO: print output to log
func (e *Executor) fetchLog(wg *sync.WaitGroup, resource *types.Resource, message *string) {
	host := resource.GetHost()
	e.content.Logger.Infof("fetch log from %s", host)
	remoteLogPath := fmt.Sprintf("%s/", path.Join(e.content.Ansible.GetDeploy(), "log"))
	localLogDir := fmt.Sprintf("log-%s", host)
	localLogPath := fmt.Sprintf("%s/", path.Join(e.logPath, localLogDir))
	args := []string{
		"-avz",
		"--port",
		strconv.Itoa(resource.GetPort()),
		fmt.Sprintf("root@%s:%s", host, remoteLogPath),
		localLogPath,
	}

	out, err := util.DoCmd(e.mgr.Config.LogPath, "rsync", args...)
	e.logger.Info(out)
	if err != nil {
		e.logger.Warnf("rsync %s failed, %v", resource.GetHost(), err)
	}
	if msg, err := e.scanLogs(host, localLogPath); err != nil {
		e.logger.Warnf("scan log failed %s failed, %v", resource.GetHost(), err)
	} else {
		*message = fmt.Sprintf("%s/%s", *message, msg)
	}
	wg.Done()
}

func (e *Executor) scanLogs(host, localLogPath string) (string, error) {
	files, err := ioutil.ReadDir(localLogPath)
	if err != nil {
		return "", errors.Trace(err)
	}
	message := ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		logFilePath := path.Join(localLogPath, file.Name())
		if logOne, err := e.scanLogFile(host, logFilePath); err != nil {
			return "", errors.Trace(err)
		} else if logOne != "" {
			message = fmt.Sprintf("%s\n%s", message, logOne)
		}
	}
	return message, nil
}

func (e *Executor) scanLogFile(host, logFilePath string) (string, error) {
	logName := logFileNameParse(logFilePath)
	if logName == "" {
		return "", nil
	}
	e.content.Logger.Infof("scan log file %s", logFilePath)
	f, err := os.Open(logFilePath)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer f.Close()

	logger := log.New(e.content.Pkg.LogEngine, log.TypeNode, e.content.Job.GetID(), e.history.GetID(),
		fmt.Sprintf("%s-%s", host, logName))
	la := logAnalyse{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		logger.Info(line)
		if strings.Contains(line, "[info]") {
			la.info++
		}
		if strings.Contains(line, "[warn]") {
			la.warn++
		}
		if strings.Contains(line, "[error]") {
			la.error++
		}
		if strings.Contains(line, "[fatal]") {
			la.fatal++
		}
	}

	//logLine := fmt.Sprintf("log file_logger: %s, info: %d, warn: %d, error: %d, fatal: %d",
	//	logFilePath, la.info, la.warn, la.error, la.fatal)
	logLine := fmt.Sprintf("INFO: %d, WARN: %d, ERROR: %d, FATAL: %d",
		la.info, la.warn, la.error, la.fatal)
	logger.Infof(logLine)

	if err := scanner.Err(); err != nil {
		return "", errors.Trace(err)
	}
	return logLine, nil
}

func logFileNameParse(filePath string) string {
	r, _ := regexp.Compile("\\/([a-zA-Z0-9_\\-\\.]+)\\.log$")
	submatch := r.FindStringSubmatch(filePath)
	if len(submatch) >= 2 {
		match := submatch[1]
		if strings.Contains(match, "tikv") {
			return match
		}
		if strings.Contains(match, "tidb") {
			return match
		}
		if strings.Contains(match, "pd") {
			return match
		}
	}
	return ""
}
