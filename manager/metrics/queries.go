package metrics

type Query struct {
	Name         string
	Stmt         string
	HigherBetter bool
}

var queries = []*Query{
	{
		Name:         "QPS",
		Stmt:         "sum(rate(tidb_server_query_total{result=\"OK\"}[30s])) by (result)",
		HigherBetter: true,
	},
	{
		Name:         "Duration99",
		Stmt:         "histogram_quantile(0.99, sum(rate(tidb_server_handle_query_duration_seconds_bucket[1m])) by (le))",
		HigherBetter: false,
	},
}
