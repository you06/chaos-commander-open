package slack

import (
	"github.com/juju/errors"
	"github.com/nlopes/slack"
	"github.com/you06/chaos-commander/config"
)

type Client struct {
	config *config.Slack
	client *slack.Client
}

func NewClient(cfg *config.Slack) (*Client, error) {
	if cfg == nil {
		return &Client{
			config: nil,
			client: nil,
		}, nil
	}

	c := Client{
		config: cfg,
		client: slack.New(cfg.Token),
	}

	if err := c.SendMessage("hello"); err != nil {
		return nil, errors.Trace(err)
	}

	return &c, nil
}

func (c *Client)SendMessage(msg string) error {
	if c.client == nil {
		return nil
	}
	if _, _, err := c.client.PostMessage(c.config.Channel,
		slack.MsgOptionText(msg, true)); err != nil {
		return errors.Trace(err)
	}
	return nil
}
