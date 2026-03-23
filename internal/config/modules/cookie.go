package modules

import "time"

// CookieConfig хранит настройки httpOnly refresh cookie.
type CookieConfig struct {
	Name     string        `yaml:"name"`
	Path     string        `yaml:"path"`
	Domain   string        `yaml:"domain"`
	Secure   *bool         `yaml:"secure"`    // nil = не задано, дефолт = true
	SameSite string        `yaml:"same_site"` // "strict", "lax", "none"
	TTL      time.Duration `yaml:"ttl"`
}

// IsSecure возвращает значение Secure-флага cookie.
// Если поле не задано — возвращает true (безопасный дефолт).
func (c CookieConfig) IsSecure() bool {
	if c.Secure == nil {
		return true
	}
	return *c.Secure
}

func (c *CookieConfig) ApplyDefaults() {
	if c.Name == "" {
		c.Name = "refresh_token"
	}
	if c.Path == "" {
		c.Path = "/api/v1/auth"
	}
	if c.SameSite == "" {
		c.SameSite = "strict"
	}
	if c.TTL == 0 {
		c.TTL = 7 * 24 * time.Hour
	}
}
