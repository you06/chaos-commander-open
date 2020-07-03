package clickhouse

import (
	"github.com/juju/errors"
	"github.com/you06/chaos-commander/pkg/log-engine/basic"
)

func (c *Conn) InsertManagerLog(line *basic.LogLine) error {
	c.Lock()
	defer c.Unlock()
	tx, err := c.conn.Begin()
	if err != nil {
		return errors.Trace(err)
	}
	stmt, err := tx.Prepare("INSERT INTO manager_logs (jobID, historyID, node, log_date, log_level, log) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Trace(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(line.JobID, line.HistoryID, line.Node, line.LogDate, line.LogLevel, line.Log); err != nil {
		return errors.Trace(err)
	}
	if err := tx.Commit(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) BatchInsertManagerLog(lines []*basic.LogLine) error {
	c.Lock()
	defer c.Unlock()
	tx, err := c.conn.Begin()
	if err != nil {
		return errors.Trace(err)
	}
	stmt, err := tx.Prepare("INSERT INTO manager_logs (jobID, historyID, node, log_date, log_level, log) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Trace(err)
	}
	defer stmt.Close()
	for _, line := range lines {
		if _, err := stmt.Exec(line.JobID, line.HistoryID, line.Node, line.LogDate, line.LogLevel, line.Log); err != nil {
			return errors.Trace(err)
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) InsertLoadLog(line *basic.LogLine) error {
	c.Lock()
	defer c.Unlock()
	tx, err := c.conn.Begin()
	if err != nil {
		return errors.Trace(err)
	}
	stmt, err := tx.Prepare("INSERT INTO load_logs (jobID, historyID, node, log_date, log_level, log) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Trace(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(line.JobID, line.HistoryID, line.Node, line.LogDate, line.LogLevel, line.Log); err != nil {
		return errors.Trace(err)
	}
	if err := tx.Commit(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) BatchInsertLoadLog(lines []*basic.LogLine) error {
	c.Lock()
	defer c.Unlock()
	tx, err := c.conn.Begin()
	if err != nil {
		return errors.Trace(err)
	}
	stmt, err := tx.Prepare("INSERT INTO load_logs (jobID, historyID, node, log_date, log_level, log) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Trace(err)
	}
	defer stmt.Close()
	for _, line := range lines {
		if _, err := stmt.Exec(line.JobID, line.HistoryID, line.Node, line.LogDate, line.LogLevel, line.Log); err != nil {
			return errors.Trace(err)
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) InsertNodeLog(line *basic.LogLine) error {
	c.Lock()
	defer c.Unlock()
	tx, err := c.conn.Begin()
	if err != nil {
		return errors.Trace(err)
	}
	stmt, err := tx.Prepare("INSERT INTO node_logs (jobID, historyID, node, log_date, log_level, log) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Trace(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(line.JobID, line.HistoryID, line.Node, line.LogDate, line.LogLevel, line.Log); err != nil {
		return errors.Trace(err)
	}
	if err := tx.Commit(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) BatchInsertNodeLog(lines []*basic.LogLine) error {
	c.Lock()
	defer c.Unlock()
	tx, err := c.conn.Begin()
	if err != nil {
		return errors.Trace(err)
	}
	stmt, err := tx.Prepare("INSERT INTO node_logs (jobID, historyID, node, log_date, log_level, log) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Trace(err)
	}
	defer stmt.Close()
	for _, line := range lines {
		if _, err := stmt.Exec(line.JobID, line.HistoryID, line.Node, line.LogDate, line.LogLevel, line.Log); err != nil {
			return errors.Trace(err)
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Trace(err)
	}

	return nil
}
