package config

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

func defaultSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:                   "smtp.protonmail.ch",
		Port:                   587,
		AuthEnabled:            true,
		FallbackTimeoutSeconds: 10,
	}
}
