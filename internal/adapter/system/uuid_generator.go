package system

import "github.com/google/uuid"

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (UUIDGenerator) New() uuid.UUID {
	return uuid.New()
}
