package modules

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	URL             string        `yaml:"url"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

func (c *DatabaseConfig) ApplyDefaults() {
	if c.URL == "" {
		c.URL = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 25
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 5
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 5 * time.Minute
	}
}

func (c DatabaseConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	return nil
}
