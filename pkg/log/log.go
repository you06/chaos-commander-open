package log

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/pkg/log-engine/basic"
	"regexp"
	"sync"
	"time"
)

type Log struct {
	sync.Mutex
	logs []*basic.LogLine
	quit chan struct{}
	Mute bool
	engine basic.Driver
	fileLogRegex *regexp.Regexp
	commonLogRegex *regexp.Regexp
	logType Type
	jobID int
	historyID int
	node string
}

func New(engine basic.Driver, logType Type, jobID, historyID int, node string) *Log {
	logger := &Log{
		logs: []*basic.LogLine{},
		quit: make(chan struct{}, 1),
		Mute: false,
		engine: engine,
		logType: logType,
		jobID: jobID,
		historyID: historyID,
		node: node,
	}
	if err := logger.initParser(); err != nil {
		log.Fatal(errors.ErrorStack(err))
	}

	go logger.watchLog()

	return logger
}

// Quit stop log
func (l *Log)Quit() {
	l.quit <- struct{}{}
}

// Write implement io.Writer interface
func (l *Log)Write(p []byte) (n int, err error) {
	l.Lock()
	defer l.Unlock()

	l.logs = append(l.logs, l.parseLine(string(p)))

	n = len(p)
	err = nil
	return
}

// Info will print a line to log
func (l *Log)Info(line string) {
	l.Lock()
	defer l.Unlock()
	l.logs = append(l.logs, l.parseLine(line))
}

// Infof will print a line to log with format
func (l *Log)Infof(line string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.logs = append(l.logs, l.parseLine(fmt.Sprintf(line, args...)))
}

// Warn will print a error log
func (l *Log)Warn(line string) {
	l.Lock()
	defer l.Unlock()
	l.logs = append(l.logs, &basic.LogLine{
		JobID:     l.jobID,
		HistoryID: l.historyID,
		Node:      l.node,
		LogDate:   time.Now(),
		LogLevel:  "[WARN]",
		Log:       line,
	})
}

// Warnf will print a error log with format
func (l *Log)Warnf(line string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.logs = append(l.logs, &basic.LogLine{
		JobID:     l.jobID,
		HistoryID: l.historyID,
		Node:      l.node,
		LogDate:   time.Now(),
		LogLevel:  "[WARN]",
		Log:       fmt.Sprintf(line, args...),
	})
}

// Error will print a error log
func (l *Log)Error(line string) {
	l.Lock()
	defer l.Unlock()
	l.logs = append(l.logs, &basic.LogLine{
		JobID:     l.jobID,
		HistoryID: l.historyID,
		Node:      l.node,
		LogDate:   time.Now(),
		LogLevel:  "[ERROR]",
		Log:       line,
	})
}

// Errorf will print a error log with format
func (l *Log)Errorf(line string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.logs = append(l.logs, &basic.LogLine{
		JobID:     l.jobID,
		HistoryID: l.historyID,
		Node:      l.node,
		LogDate:   time.Now(),
		LogLevel:  "[ERROR]",
		Log:       fmt.Sprintf(line, args...),
	})
}

func (l *Log)watchLog() {
	ticker := time.NewTicker(1 * time.Second)
	forceFlushTicker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <- ticker.C:
			if len(l.logs) >= BATCH_SIZE {
				if err := l.flush(); err != nil {
					log.Error(errors.ErrorStack(err))
				}
			}
		case <- forceFlushTicker.C:
			if err := l.flush(); err != nil {
				log.Error(errors.ErrorStack(err))
			}
		case <- l.quit:
			ticker.Stop()
			if err := l.flush(); err != nil {
				log.Error(err)
			}
			return
		}
	}
}

func (l *Log)flush() error {
	l.Lock()
	if len(l.logs) == 0 {
		l.Unlock()
		if err := l.engine.Ping(); err != nil {
			log.Errorf("engine ping failed %v", err)
		}
		return nil
	}
	flushLogs := l.logs
	l.logs = []*basic.LogLine{}
	l.Unlock()

	var err error
	switch l.logType {
	case TypeManager:
		err = l.engine.BatchInsertManagerLog(flushLogs)
	case TypeLoad:
		err = l.engine.BatchInsertLoadLog(flushLogs)
	case TypeNode:
		err = l.engine.BatchInsertNodeLog(flushLogs)
	}
	if err != nil {
		log.Error(errors.ErrorStack(err))
	}
	//f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	return errors.Trace(err)
	//}
	//defer func () {
	//	if err := f.Close(); err != nil {
	//		log.Error(err)
	//	}
	//}()
	//
	//if _, err = f.WriteString(strings.Join(flushLogs, "\n")); err != nil {
	//	return errors.Trace(err)
	//}

	return nil
}

func getDateStr() string {
	now := time.Now()
	return now.Format("2006/01/02 15:04:05")
}
