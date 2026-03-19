package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/dto"
	domainservice "github.com/AlbinaKonovalova/auth-service/internal/domain/service"
)

// CreateUserInput содержит данные для создания нового пользователя.
type CreateUserInput struct {
	Email    string
	Password string
	Roles    []string
}

// UserUseCase — входной контракт для users admin API.
type UserUseCase interface {
	// CreateUser создаёт нового пользователя с опциональными начальными ролями.
	CreateUser(ctx context.Context, in CreateUserInput) (domainservice.AdminUserView, error)

	// GetUser возвращает пользователя по его ID вместе со списком ролей.
	GetUser(ctx context.Context, id uuid.UUID) (domainservice.AdminUserView, error)

	// ListUsers возвращает постраничный список пользователей с применением фильтров.
	ListUsers(ctx context.Context, filters dto.UserListFilters) (domainservice.UserList, error)

	// ActivateUser активирует учётную запись пользователя.
	ActivateUser(ctx context.Context, id uuid.UUID) error

	// DeactivateUser деактивирует учётную запись пользователя и отзывает все его refresh-сессии.
	DeactivateUser(ctx context.Context, id uuid.UUID) error
}
