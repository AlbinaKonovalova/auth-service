package modules

import (
	"fmt"
	"time"
)

// PasswordResetConfig хранит настройки сброса пароля.
type PasswordResetConfig struct {
	TTL     time.Duration `yaml:"ttl"`
	BaseURL string        `yaml:"base_url"`
}

func (c *PasswordResetConfig) ApplyDefaults() {
	if c.TTL == 0 {
		c.TTL = 30 * time.Minute
	}
	if c.BaseURL == "" {
		c.BaseURL = "http://localhost:3000/reset-password"
	}
}

func (c PasswordResetConfig) Validate() error {
	if c.TTL <= 0 {
		return fmt.Errorf("password_reset.ttl must be greater than 0")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("password_reset.base_url is required")
	}
	return nil
}
