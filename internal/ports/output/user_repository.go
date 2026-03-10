package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/entity"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email value.Email) (*entity.User, error)
}
