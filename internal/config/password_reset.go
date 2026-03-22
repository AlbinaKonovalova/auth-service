package config

import "time"

type PasswordResetConfig struct {
	TTL     time.Duration `yaml:"ttl"`
	BaseURL string        `yaml:"base_url"`
}

func defaultPasswordResetConfig() PasswordResetConfig {
	return PasswordResetConfig{
		TTL:     30 * time.Minute,
		BaseURL: "http://localhost:3000/reset-password",
	}
}
