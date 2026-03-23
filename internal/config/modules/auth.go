package modules

import (
	"fmt"
	"time"
)

// AuthConfig хранит настройки JWT и refresh token.
type AuthConfig struct {
	JWTSecret         string        `yaml:"jwt_secret"`
	AccessTokenTTL    time.Duration `yaml:"access_token_ttl"`
	RefreshTokenBytes int           `yaml:"refresh_token_bytes"`
}

func (c *AuthConfig) ApplyDefaults() {
	if c.AccessTokenTTL == 0 {
		c.AccessTokenTTL = 15 * time.Minute
	}
	if c.RefreshTokenBytes == 0 {
		c.RefreshTokenBytes = 32
	}
}

func (c AuthConfig) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}
	return nil
}
