package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/AlbinaKonovalova/auth-service/internal/config/modules"
)

type Config struct {
	Server        modules.ServerConfig        `yaml:"server"`
	Database      modules.DatabaseConfig      `yaml:"database"`
	Log           modules.LogConfig           `yaml:"log"`
	Cookie        modules.CookieConfig        `yaml:"cookie"`
	Auth          modules.AuthConfig          `yaml:"auth"`
	Argon2        modules.Argon2Config        `yaml:"argon2"`
	CORS          modules.CORSConfig          `yaml:"cors"`
	SMTP          modules.SMTPConfig          `yaml:"smtp"`
	PasswordReset modules.PasswordResetConfig `yaml:"password_reset"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}

	if err := cfg.loadFromFile(path); err != nil {
		return nil, fmt.Errorf("failed to load config from file: %w", err)
	}

	cfg.setDefaults()

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) loadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

func (c *Config) setDefaults() {
	c.Server.ApplyDefaults()
	c.Database.ApplyDefaults()
	c.Log.ApplyDefaults()
	c.Auth.ApplyDefaults()
	c.Cookie.ApplyDefaults()
	c.Argon2.ApplyDefaults()
	c.CORS.ApplyDefaults()
	c.SMTP.ApplyDefaults()
	c.PasswordReset.ApplyDefaults()
}

func (c *Config) validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}
	if err := c.Database.Validate(); err != nil {
		return err
	}
	if err := c.Auth.Validate(); err != nil {
		return err
	}
	if err := c.SMTP.Validate(); err != nil {
		return err
	}
	if err := c.PasswordReset.Validate(); err != nil {
		return err
	}
	return nil
}
