package config

// Argon2Config содержит параметры argon2id для хеширования паролей.
type Argon2Config struct {
	Memory      uint32 `yaml:"memory"`
	Iterations  uint32 `yaml:"iterations"`
	Parallelism uint8  `yaml:"parallelism"`
	SaltLength  uint32 `yaml:"salt_length"`
	KeyLength   uint32 `yaml:"key_length"`
}

func defaultArgon2Config() Argon2Config {
	return Argon2Config{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}
