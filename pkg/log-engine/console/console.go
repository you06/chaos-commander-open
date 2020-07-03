package console

import (
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/pkg/log-engine/basic"
)

type Console struct{}

func New() *Console {
	return &Console{}
}

func (c *Console) Ping() error {
	return nil
}

func (c *Console) InsertManagerLog(line *basic.LogLine) error {
	log.Infof("%s %s %s", "manager", line.LogLevel, line.Log)
	return nil
}

func (c *Console) InsertLoadLog(line *basic.LogLine) error {
	log.Infof("%s %s %s", "manager", line.LogLevel, line.Log)
	return nil
}

func (c *Console) InsertNodeLog(line *basic.LogLine) error {
	log.Infof("%s %s %s", "manager", line.LogLevel, line.Log)
	return nil
}

func (c *Console) BatchInsertManagerLog(lines []*basic.LogLine) error {
	for _, line := range lines {
		log.Infof("%s %s %s", "manager", line.LogLevel, line.Log)
	}
	return nil
}

func (c *Console) BatchInsertLoadLog(lines []*basic.LogLine) error {
	for _, line := range lines {
		log.Infof("%s %s %s", "manager", line.LogLevel, line.Log)
	}
	return nil
}

func (c *Console) BatchInsertNodeLog(lines []*basic.LogLine) error {
	for _, line := range lines {
		log.Infof("%s %s %s", "manager", line.LogLevel, line.Log)
	}
	return nil
}
