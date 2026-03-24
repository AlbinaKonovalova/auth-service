package modules

import "fmt"

type SMTPConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	AuthEnabled            bool   `yaml:"auth_enabled"`
	Username               string `yaml:"username"`
	Password               string `yaml:"password"`
	From                   string `yaml:"from"`
	FallbackTimeoutSeconds int    `yaml:"fallback_timeout_seconds"`
	SkipTLSVerify          bool   `yaml:"skip_tls_verify"`
}

func (c *SMTPConfig) ApplyDefaults() {
	if c.Host == "" {
		c.Host = "smtp.protonmail.ch"
	}
	if c.Port == 0 {
		c.Port = 587
	}
	if c.FallbackTimeoutSeconds == 0 {
		c.FallbackTimeoutSeconds = 10
	}
}

func (c SMTPConfig) Validate() error {
	if c.From == "" {
		return fmt.Errorf("smtp.from is required")
	}
	if c.Host == "" {
		return fmt.Errorf("smtp.host is required")
	}
	if c.Port <= 0 {
		return fmt.Errorf("smtp.port must be greater than 0")
	}
	if c.FallbackTimeoutSeconds <= 0 {
		return fmt.Errorf("smtp.fallback_timeout_seconds must be greater than 0")
	}
	if c.AuthEnabled {
		if c.Username == "" || c.Password == "" {
			return fmt.Errorf("smtp.username and smtp.password are required when smtp.auth_enabled=true")
		}
	} else {
		if c.Username != "" || c.Password != "" {
			return fmt.Errorf("smtp.username and smtp.password must be empty when smtp.auth_enabled=false")
		}
	}
	return nil
}
