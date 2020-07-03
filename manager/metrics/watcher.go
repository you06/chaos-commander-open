package metrics

import (
	"time"

	"github.com/ngaut/log"
)

func (m *Metrics) StartWatch() {
	log.Info("start watch")
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				m.checkOnce()
			case <-m.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *Metrics) StopWatch() {
	m.quit <- struct{}{}
}

func (m *Metrics) checkOnce() {
	for _, query := range queries {
		go func(query *Query) {
			v, err := m.QueryProm(query.Stmt, time.Now())
			if err != nil {
				log.Infof("error fetch metrics %v", err)
				return
			}

			newResult := QueryResult{
				Name:  query.Name,
				Stmt:  query.Stmt,
				Value: *v,
				Time:  time.Now(),
			}
			if oldResult, exist := m.LastResult[query.Name]; exist {
				log.Infof("%s, value %f, last value %f", query.Name, newResult.Value, oldResult.Value)
				// compare
				increase := (newResult.Value - oldResult.Value) > 0
				if increase != query.HigherBetter {
					// calc the change rate
					rate := 100 * (newResult.Value - oldResult.Value) / (oldResult.Value + 0.000001)
					action := "increase"
					if rate < 0 {
						rate = -rate
						action = "drop"
					}
					if rate > 0 {
						log.Infof("%s %s %f", query.Name, action, rate)
					}
				}
			}
			m.LastResult[query.Name] = &newResult
		}(query)
	}
}
