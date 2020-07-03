package clickhouse

import (
	"database/sql"
	"fmt"
	"sync"

	clickhouse "github.com/ClickHouse/clickhouse-go"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/chaos-commander/config"
)

type Conn struct {
	dsn  string
	cfg  *config.Clickhouse
	conn *sql.DB
	sync.Mutex
}

func New(cfg *config.Clickhouse) (*Conn, error) {
	dsn := fmt.Sprintf("tcp://%s:%d?&username=%s&password=%s&database=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
	log.Infof("dsn: %s", dsn)
	conn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, errors.Trace(err)
	}

	c := &Conn{
		dsn:  dsn,
		cfg:  cfg,
		conn: conn,
	}

	if err := c.init(); err != nil {
		return nil, errors.Trace(err)
	}

	return c, nil
}

func (c *Conn) init() error {
	if err := c.Ping(); err != nil {
		return errors.Trace(err)
	}

	if err := c.mustExec(`
		CREATE TABLE IF NOT EXISTS manager_logs (
			date         Date MATERIALIZED toDate(log_date),
			jobID        Int32,
			historyID    Int32,
			node         String,
			log_date     DateTime,
			log_level    String,
			log          String
		) engine=MergeTree(date, (jobID, historyID, log_date), 32768)
	`); err != nil {
		return errors.Trace(err)
	}

	if err := c.mustExec(`
		CREATE TABLE IF NOT EXISTS load_logs (
			date         Date MATERIALIZED toDate(log_date),
			jobID        Int32,
			historyID    Int32,
			node         String,
			log_date     DateTime,
			log_level    String,
			log          String
		) engine=MergeTree(date, (jobID, historyID, log_date), 32768)
	`); err != nil {
		return errors.Trace(err)
	}

	if err := c.mustExec(`
		CREATE TABLE IF NOT EXISTS node_logs (
			date         Date MATERIALIZED toDate(log_date),
			jobID        Int32,
			historyID    Int32,
			node         String,
			log_date     DateTime,
			log_level    String,
			log          String
		) engine=MergeTree(date, (jobID, historyID, log_date), 32768)
	`); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) Ping() error {
	if err := c.conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return errors.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			return errors.Trace(err)
		}
	}
	return nil
}

func (c *Conn) mustExec(sql string) error {
	_, err := c.conn.Exec(sql)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (c *Conn) reConn() {
	c.Lock()
	defer c.Unlock()
	if err := c.conn.Close(); err != nil {
		log.Errorf("conn close error, %v", err)
	}
	conn, err := sql.Open("clickhouse", c.dsn)
	if err != nil {
		log.Errorf("conn open error, %v", err)
	} else {
		c.conn = conn
	}
}
