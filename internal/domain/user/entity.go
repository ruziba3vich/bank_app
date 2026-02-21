package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	FullName     string
	Role         string
	Login        string
	PasswordHash string
	Status       bool
	CreatedAt    time.Time
}

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	DeviceID         string
	RefreshTokenHash string
	LastLoginAt      time.Time
	CreatedAt        time.Time
}
