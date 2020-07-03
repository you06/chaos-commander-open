package core

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/config"
	"github.com/you06/chaos-commander/pkg"
	"github.com/you06/chaos-commander/pkg/storage/basic"
)

// Manager represent schrodinger manager
type Manager struct {
	Config  *config.Config
	Pkg     *pkg.Pkg
	storage basic.Driver
}

// New init manager
func New(cfg *config.Config) (*Manager, error) {
	p, err := pkg.New(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Manager{
		Config:  cfg,
		Pkg:     p,
		storage: p.Storage,
	}, nil
}
