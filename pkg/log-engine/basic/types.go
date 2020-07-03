package basic

import (
	"fmt"
	"time"
)

type LogLine struct {
	JobID int
	HistoryID int
	Node string
	LogDate time.Time
	LogLevel string
	Log string
}

func (l *LogLine)String() string {
	return fmt.Sprintf("%s %s %s",
		l.LogDate.Format("2006/01/02 15:04:05"), l.LogLevel, l.Log)
}
