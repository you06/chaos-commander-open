package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLogRegex(t *testing.T) {
	assert.Equal(t, logFileNameParse("/home/pingcap/deploy-partition/log/tidb_stderr.log"),
		"tidb_stderr", "tidb_stderr log parse")
	assert.Equal(t, logFileNameParse("/home/tidb/log/disk-write-2019-09-29T16:04:00/172.16.5.16/tidb.log"),
		"tidb", "tidb log parse")
	assert.Equal(t, logFileNameParse("/home/tidb/log/disk-write-2019-09-29T16:04:00/172.16.5.16/tikv.log"),
		"tikv", "tikv log parse")
	assert.Equal(t, logFileNameParse("/home/tidb/log/disk-write-2019-09-29T16:04:00/172.16.5.16/pd.log"),
		"pd", "pd log parse")
	assert.Equal(t, logFileNameParse("/home/tidb/log/disk-write-2019-09-29T16:04:00/172.16.5.16/node_exporter.log"),
		"", "node_exporter log parse failed")
}
