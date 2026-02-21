package dto

import (
	"time"

	"github.com/google/uuid"

	domain "github.com/prodonik/bank_app/internal/domain/user"
)

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceID     string `json:"device_id" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	Login     string    `json:"login"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Role:      u.Role,
		Login:     u.Login,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}

func NewAuthResponse(accessToken, refreshToken string, u *domain.User) AuthResponse {
	return AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         NewUserResponse(u),
	}
}
