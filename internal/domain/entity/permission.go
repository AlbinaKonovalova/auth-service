package entity

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID
	Code        string
	Description string
}
