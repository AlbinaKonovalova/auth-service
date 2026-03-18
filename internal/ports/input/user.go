package input

import (
	"context"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	"github.com/google/uuid"
)

// CreateUserInput содержит данные для создания нового пользователя.
type CreateUserInput struct {
	Email    string
	Password string

	// Roles — список кодов ролей, назначаемых при создании (может быть пустым).
	Roles []string
}

// CreateUserResult — результат создания пользователя.
type CreateUserResult struct {
	ID       uuid.UUID
	Email    string
	IsActive bool
	Roles    []string
}

// GetUserResult — результат получения пользователя по ID.
type GetUserResult struct {
	ID       uuid.UUID
	Email    string
	IsActive bool
	Roles    []string
}

// ListUsersResult — результат постраничного списка пользователей.
type ListUsersResult struct {
	Users   []GetUserResult
	Total   int
	Page    int
	PerPage int
}

// UserUseCase — входной контракт для users admin API.
type UserUseCase interface {
	// CreateUser создаёт нового пользователя с опциональными начальными ролями.
	CreateUser(ctx context.Context, in CreateUserInput) (CreateUserResult, error)

	// GetUser возвращает пользователя по его ID вместе со списком ролей.
	GetUser(ctx context.Context, id uuid.UUID) (GetUserResult, error)

	// ListUsers возвращает постраничный список пользователей с применением фильтров.
	ListUsers(ctx context.Context, filters dto.UserListFilters) (ListUsersResult, error)

	// ActivateUser активирует учётную запись пользователя.
	ActivateUser(ctx context.Context, id uuid.UUID) error

	// DeactivateUser деактивирует учётную запись пользователя и отзывает все его refresh-сессии.
	DeactivateUser(ctx context.Context, id uuid.UUID) error
}
