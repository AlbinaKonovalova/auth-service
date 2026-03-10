package auth

import (
	"time"

	"github.com/google/uuid"
)

// Clock реализует output.Clock.
type Clock struct{}

func NewClock() *Clock {
	return &Clock{}
}

func (Clock) Now() time.Time {
	return time.Now()
}

// UUIDGenerator реализует output.UUIDGenerator.
type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (UUIDGenerator) New() uuid.UUID {
	return uuid.New()
}
