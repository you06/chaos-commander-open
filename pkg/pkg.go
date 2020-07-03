package pkg

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/config"
	logEngineBasic "github.com/you06/chaos-commander/pkg/log-engine/basic"
	"github.com/you06/chaos-commander/pkg/log-engine/clickhouse"
	"github.com/you06/chaos-commander/pkg/log-engine/console"
	"github.com/you06/chaos-commander/pkg/slack"
	"github.com/you06/chaos-commander/pkg/storage"
	"github.com/you06/chaos-commander/pkg/storage/basic"
)

type Pkg struct {
	LogEngine logEngineBasic.Driver
	Storage   basic.Driver
	Slack     *slack.Client
}

func New(cfg *config.Config) (*Pkg, error) {
	store, err := storage.New(cfg.Database)
	if err != nil {
		return nil, errors.Trace(err)
	}

	slackClient, err := slack.NewClient(cfg.Slack)
	if err != nil {
		return nil, errors.Trace(err)
	}

	pkg := Pkg{
		LogEngine: nil,
		Storage:   store,
		Slack:     slackClient,
	}

	switch cfg.LogEngine {
	case "clickhouse":
		chouse, err := clickhouse.New(cfg.Clickhouse)
		if err != nil {
			return nil, errors.Trace(err)
		}
		pkg.LogEngine = chouse
	case "console":
		pkg.LogEngine = console.New()
	default:
		pkg.LogEngine = console.New()
	}

	return &pkg, nil
}
