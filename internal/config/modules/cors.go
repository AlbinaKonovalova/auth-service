package modules

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

func (c *CORSConfig) ApplyDefaults() {
	if c.AllowedOrigins == nil {
		c.AllowedOrigins = []string{}
	}
}
