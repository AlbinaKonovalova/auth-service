package config

import "time"

// CookieConfig содержит настройки refresh token cookie.
type CookieConfig struct {
	Name     string        `yaml:"name"`
	Path     string        `yaml:"path"`
	Domain   string        `yaml:"domain"`
	Secure   *bool         `yaml:"secure"`    // указатель: nil = не задано, дефолт = true
	SameSite string        `yaml:"same_site"` // "strict", "lax", "none"
	TTL      time.Duration `yaml:"ttl"`
}

// IsSecure возвращает значение Secure с дефолтом true.
func (c CookieConfig) IsSecure() bool {
	if c.Secure == nil {
		return true
	}
	return *c.Secure
}

func defaultCookieConfig() CookieConfig {
	secure := true
	return CookieConfig{
		Name:     "refresh_token",
		Path:     "/auth",
		Domain:   "",
		Secure:   &secure,
		SameSite: "strict",
		TTL:      7 * 24 * time.Hour,
	}
}
