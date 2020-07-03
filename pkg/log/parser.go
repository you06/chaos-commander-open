package log

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/log-engine/basic"
	"regexp"
	"strings"
	"time"
)

func (l *Log)initParser() error {
	fileLogRegex, err := regexp.Compile(fileLogRegexString)
	if err != nil {
		return errors.Trace(err)
	}
	commonLogRegex, err := regexp.Compile(commonLogRegexString)
	if err != nil {
		return errors.Trace(err)
	}
	l.fileLogRegex = fileLogRegex
	l.commonLogRegex = commonLogRegex

	return nil
}

func (l *Log)parseLine(line string) *basic.LogLine {
	// match tikv/tidb/pd log
	// [date] [level] [log]...
	if logLine, match := l.matchFileLog(line); match {
		return logLine
	}
	// match common
	// date log
	if logLine, match := l.matchCommonLog(line); match {
		return logLine
	}

	return &basic.LogLine{
		JobID: l.jobID,
		HistoryID: l.historyID,
		Node: l.node,
		LogDate: time.Now(),
		LogLevel: "[INFO]",
		Log: line,
	}
}

func (l *Log)matchFileLog(line string) (*basic.LogLine, bool) {
	match := l.fileLogRegex.FindStringSubmatch(line)
	if len(match) == 4 {
		t, err := time.Parse("2006/01/02 15:04:05.000 -07:00", match[1])
		if err != nil {
			return nil, false
		}

		return &basic.LogLine{
			JobID: l.jobID,
			HistoryID: l.historyID,
			Node: l.node,
			LogLevel: strings.ToUpper(match[2]),
			LogDate: t,
			Log: match[3],
		}, true
	}
	return nil, false
}

func (l *Log)matchCommonLog(line string) (*basic.LogLine, bool) {
	match := l.commonLogRegex.FindStringSubmatch(line)
	if len(match) == 3 {
		t, err := time.Parse("2006/01/02 15:04:05", match[1])
		if err != nil {
			return nil, false
		}

		return &basic.LogLine{
			JobID: l.jobID,
			HistoryID: l.historyID,
			Node: l.node,
			LogLevel: strings.ToUpper(match[2]),
			LogDate: t,
			Log: match[0],
		}, true
	}
	return nil, false
}
