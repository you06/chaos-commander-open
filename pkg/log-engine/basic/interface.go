package basic

// Driver defines log engine model
type Driver interface {
	Ping() error
	InsertManagerLog(line *LogLine) error
	InsertLoadLog(line *LogLine) error
	InsertNodeLog(line *LogLine) error
	BatchInsertManagerLog(lines []*LogLine) error
	BatchInsertLoadLog(lines []*LogLine) error
	BatchInsertNodeLog(lines []*LogLine) error
}
