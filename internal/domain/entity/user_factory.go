package entity

import (
	"time"

	domain "github.com/AlbinaKonovalova/auth-service/internal/domain"
	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
	"github.com/google/uuid"
)

type NewUserAggregate struct {
	User          User
	UserRoles     []UserRole
	AssignedRoles []string
}
 fun

