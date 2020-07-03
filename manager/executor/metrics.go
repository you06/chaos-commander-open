package executor

func (e *Executor) metricsWatcherStart() {
	e.metrics.StartWatch()
}

func (e *Executor) metricsWatcherStop() {
	e.metrics.StopWatch()
}
