package entity

import (
	"strings"
	"time"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, email value.Email, passwordHash string, now time.Time) (User, error) {
	if id == uuid.Nil {
		return User{}, domain.ErrInvalidUserID
	}
	if strings.TrimSpace(passwordHash) == "" {
		return User{}, domain.ErrInvalidPasswordHash
	}
	if now.IsZero() {
		return User{}, domain.ErrInvalidUserID
	}

	return User{
		ID:           id,
		Email:        email.String(),
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u User) EnsureActive() error {
	if !u.IsActive {
		return domain.ErrUserInactive
	}

	return nil
}
