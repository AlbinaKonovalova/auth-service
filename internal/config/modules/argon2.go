package modules

// Argon2Config хранит параметры алгоритма argon2id для хеширования паролей.
type Argon2Config struct {
	Memory      uint32 `yaml:"memory"`
	Iterations  uint32 `yaml:"iterations"`
	Parallelism uint8  `yaml:"parallelism"`
	SaltLength  uint32 `yaml:"salt_length"`
	KeyLength   uint32 `yaml:"key_length"`
}

func (c *Argon2Config) ApplyDefaults() {
	if c.Memory == 0 {
		c.Memory = 64 * 1024
	}
	if c.Iterations == 0 {
		c.Iterations = 3
	}
	if c.Parallelism == 0 {
		c.Parallelism = 2
	}
	if c.SaltLength == 0 {
		c.SaltLength = 16
	}
	if c.KeyLength == 0 {
		c.KeyLength = 32
	}
}
