package config

// CORSConfig содержит настройки CORS политики.
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

func defaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{},
	}
}
