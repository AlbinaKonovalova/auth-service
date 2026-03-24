package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

type CreateUserInput struct {
	Email    string
	Password string
	Roles    []string
}

type UserUseCase interface {
	CreateUser(ctx context.Context, in CreateUserInput) (domainservice.AdminUserView, error)

	GetUser(ctx context.Context, id uuid.UUID) (domainservice.AdminUserView, error)

	ListUsers(ctx context.Context, filters dto.UserListFilters) (domainservice.UserList, error)

	ActivateUser(ctx context.Context, id uuid.UUID) error

	DeactivateUser(ctx context.Context, id uuid.UUID) error
}
