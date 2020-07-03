package file_logger

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/log-engine/basic"
)

func (l *Logger) InsertJobLog(line *basic.LogLine) error {
	if err := l.writeLine(line.String()); err != nil {
		return errors.Trace(err)
	}
	return nil
}
