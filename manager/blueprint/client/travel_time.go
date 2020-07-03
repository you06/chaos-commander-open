package client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/juju/errors"
	"github.com/you06/chaos-commander/manager/executor/content"
	"github.com/you06/chaos-commander/pkg/types"
)

type TravelTime struct {
	resource *types.Resource
}

func newTravelTime(resource *types.Resource) *TravelTime {
	client := TravelTime{
		resource: resource,
	}
	return &client
}

func (t *TravelTime) travelTimeRandom() error {
	client, err := newSSHClient(t.resource)
	if err != nil {
		return errors.Trace(err)
	}

	err = client.Run(fmt.Sprintf("date --set=\"%s\"", randomTime()))
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (t *TravelTime) travelTimeClear() error {
	client, err := newSSHClient(t.resource)
	if err != nil {
		return errors.Trace(err)
	}

	if err := client.Run("systemctl stop ntp"); err != nil {
		return errors.Trace(err)
	}
	if err := client.Run("ntpdate -u pool.ntp.org"); err != nil {
		return errors.Trace(err)
	}
	if err := client.Run("systemctl start ntp"); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func randomTime() string {
	rand.Seed(time.Now().UnixNano())

	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min

	t := time.Unix(sec, 0)
	return t.Format("2 Jan 2006 15:04:05")
}

func TravelTimeDo(ctn *content.Content, resource *types.Resource, args ...interface{}) error {
	s := newTravelTime(resource)
	return errors.Trace(s.travelTimeRandom())
}

func TravelTimeClear(ctn *content.Content, resource *types.Resource, args ...interface{}) error {
	s := newTravelTime(resource)
	return errors.Trace(s.travelTimeClear())
}
