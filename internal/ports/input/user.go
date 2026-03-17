package input

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/google/uuid"
)

type CreateUserInput struct {
	Email    string
	Password string
	Roles    []string
}
type CreateUserResult struct {
	ID       uuid.UUID
	Email    string
	IsActive bool
	Roles    []string
}

type GetUserResult struct {
	ID       uuid.UUID
	Email    string
	IsActive bool
	Roles    []string
}

type ListUsersResult struct {
	Users   []GetUserResult
	Total   int
	Page    int
	PerPage int
}

type UserUseCase interface {
	CreateUser(ctx context.Context, in CreateUserInput) (CreateUserResult, error)
	GetUser(ctx context.Context, id uuid.UUID) (GetUserResult, error)
	ListUsers(ctx context.Context, filters dto.UserListFilters) (ListUsersResult, error)
	ActivateUser(ctx context.Context, id uuid.UUID) error
	DeactivateUser(ctx context.Context, id uuid.UUID) error
}
