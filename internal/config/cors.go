package config

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

func defaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{},
	}
}
