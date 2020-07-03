package storage

import (
	"github.com/you06/chaos-commander/config"
	"github.com/you06/chaos-commander/pkg/storage/basic"
	"github.com/you06/chaos-commander/pkg/storage/mysql"

	"github.com/juju/errors"
)

func New(config *config.Database) (basic.Driver, error) {
	driver, err := mysql.Connect(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return driver, nil
}
