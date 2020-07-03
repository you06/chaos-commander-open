package metrics

import (
	"context"
	"errors"
	"time"

	"github.com/ngaut/log"
	"github.com/prometheus/client_golang/api"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	pm "github.com/prometheus/common/model"
)

type QueryResult struct {
	Name  string
	Stmt  string
	Value float64
	Time  time.Time
}

type WatchMetric struct {
	Name     string
	Stmt     string
	Start    time.Time
	Duration time.Duration
}

type Metrics struct {
	addr       string
	quit       chan struct{}
	LastResult map[string]*QueryResult
	Watches    map[string]*WatchMetric
}

func New(addr string) *Metrics {
	log.Infof("prom addr %s", addr)
	return &Metrics{
		addr:       addr,
		quit:       make(chan struct{}, 1),
		LastResult: make(map[string]*QueryResult),
		Watches:    make(map[string]*WatchMetric),
	}
}

// QueryProm perform simple query
func (m *Metrics) QueryProm(query string, t time.Time) (*float64, error) {
	client, err := api.NewClient(api.Config{Address: m.addr})
	if err != nil {
		return nil, err
	}
	promAPI := prom.NewAPI(client)
	v, err := promAPI.Query(context.Background(), query, t)
	if err != nil {
		return nil, err
	}

	vec, ok := v.(pm.Vector)
	if !ok {
		return nil, errors.New("query prometheus: result type mismatch")
	}

	if len(vec) == 0 {
		return nil, errors.New("metric not found")
	}

	//log.Info(vec)
	value := float64(vec[0].Value)
	return &value, nil
}

// QueryPromRange perform query with duration range
func (m *Metrics) QueryPromRange(query string, start, end time.Time, step time.Duration) (FloatArray, error) {
	values := FloatArray{}
	client, err := api.NewClient(api.Config{Address: m.addr})
	if err != nil {
		return values, err
	}
	promAPI := prom.NewAPI(client)
	v, err := promAPI.QueryRange(context.Background(), query, prom.Range{start, end, step})
	if err != nil {
		return values, err
	}

	mat, ok := v.(pm.Matrix)
	if !ok {
		return values, errors.New("query prometheus: result type mismatch")
	}

	if len(mat) == 0 {
		return values, nil
	}

	for _, v := range mat[0].Values {
		values = append(values, float64(v.Value))
	}
	return values, nil
}
