package config

import "time"

type AuthConfig struct {
	JWTSecret         string        `yaml:"jwt_secret"`
	AccessTokenTTL    time.Duration `yaml:"access_token_ttl"`
	RefreshTokenBytes int           `yaml:"refresh_token_bytes"`
}

func defaultAuthConfig() AuthConfig {
	return AuthConfig{
		AccessTokenTTL:    15 * time.Minute,
		RefreshTokenBytes: 32,
	}
}
