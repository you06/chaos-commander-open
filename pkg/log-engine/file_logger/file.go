package file_logger

import (
	"os"

	"github.com/juju/errors"
	"github.com/ngaut/log"
)

type Logger struct {
	logPath string
}

func New(logPath string) (*Logger, error) {
	logger := Logger{
		logPath: logPath,
	}
	if err := logger.init(); err != nil {
		return nil, errors.Trace(err)
	}

	return &logger, nil
}

func (l *Logger) init() error {
	if err := l.writeLine("start file_logger log"); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (l *Logger) writeLine(line string) error {
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = f.WriteString(line); err != nil {
		return errors.Trace(err)
	}

	return nil
}
