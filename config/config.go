package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Config contains configuration options.
type Config struct {
	Host       string      `toml:"host"`
	Port       int         `toml:"port"`
	Root       string      `toml:"root"`
	LogPath    string      `toml:"logPath"`
	LogLevel   string      `toml:"logLevel"`
	Database   *Database   `toml:"database"`
	Slack      *Slack      `toml:"slack"`
	LogEngine  string      `toml:"logEngine"`
	Clickhouse *Clickhouse `toml:"clickhouse"`
}

// Database defines db configuration
type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

// Slack defines slack
type Slack struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

// Clickhouse defines click house configuration
type Clickhouse struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

var globalConf = Config{
	Host:     "0.0.0.0",
	Port:     50000,
	Root:     "./web/dist",
	LogPath:  "/tmp/chaos-commander",
	LogLevel: "debug",
	Database: &Database{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "",
		Database: "chaos_commander",
	},
	Slack:     nil,
	LogEngine: "console",
	Clickhouse: &Clickhouse{
		Host:     "127.0.0.1",
		Port:     9000,
		User:     "default",
		Password: "",
		Database: "chaos_commander",
	},
}

// GetGlobalConfig returns the global configuration for this server.
func GetGlobalConfig() *Config {
	return &globalConf
}

// Load loads config options from a toml file_logger.
func (c *Config) Load(confFile string) error {
	_, err := toml.DecodeFile(confFile, c)
	return errors.Trace(err)
}

// Init do some prepare works
func (c *Config) Init() error {
	if _, err := os.Stat(c.LogPath); os.IsNotExist(err) {
		if er := os.MkdirAll(c.LogPath, 0777); er != nil {
			return errors.Trace(er)
		}
	}
	return nil
}
