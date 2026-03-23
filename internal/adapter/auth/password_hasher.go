package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/AlbinaKonovalova/auth-service/internal/config/modules"
)

// decodedArgon2Params используется только при декодировании хеша из БД.
type decodedArgon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type PasswordHasher struct {
	cfg modules.Argon2Config
}

func NewPasswordHasher(cfg modules.Argon2Config) *PasswordHasher {
	return &PasswordHasher{cfg: cfg}
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, h.cfg.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.cfg.Iterations,
		h.cfg.Memory,
		h.cfg.Parallelism,
		h.cfg.KeyLength,
	)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.cfg.Memory,
		h.cfg.Iterations,
		h.cfg.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

func (h *PasswordHasher) Verify(password, encoded string) (bool, error) {
	params, salt, hash, err := decodeHash(encoded)
	if err != nil {
		return false, fmt.Errorf("decode hash: %w", err)
	}

	comparison := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	if subtle.ConstantTimeCompare(hash, comparison) != 1 {
		return false, nil
	}

	return true, nil
}

func decodeHash(encoded string) (decodedArgon2Params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return decodedArgon2Params{}, nil, nil, errors.New("invalid hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return decodedArgon2Params{}, nil, nil, errors.New("invalid hash version")
	}
	if version != argon2.Version {
		return decodedArgon2Params{}, nil, nil, errors.New("incompatible argon2 version")
	}

	var params decodedArgon2Params
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism); err != nil {
		return decodedArgon2Params{}, nil, nil, errors.New("invalid hash params")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return decodedArgon2Params{}, nil, nil, errors.New("invalid hash salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return decodedArgon2Params{}, nil, nil, errors.New("invalid hash key")
	}

	params.KeyLength = uint32(len(hash))
	params.SaltLength = uint32(len(salt))

	return params, salt, hash, nil
}
