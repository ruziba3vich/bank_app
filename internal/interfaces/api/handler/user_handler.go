package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appuser "github.com/prodonik/bank_app/internal/application/user"
	domain "github.com/prodonik/bank_app/internal/domain/user"
	"github.com/prodonik/bank_app/internal/interfaces/api/dto"
)

type UserHandler struct {
	userService *appuser.Service
}

func NewUserHandler(userService *appuser.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), appuser.RegisterInput{
		FullName: req.FullName,
		Login:    req.Login,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.NewUserResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	output, err := h.userService.Login(c.Request.Context(), appuser.LoginInput{
		Login:    req.Login,
		Password: req.Password,
		DeviceID: req.DeviceID,
	})
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewAuthResponse(output.AccessToken, output.RefreshToken, output.User))
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	output, err := h.userService.Refresh(c.Request.Context(), appuser.RefreshInput{
		RefreshToken: req.RefreshToken,
		DeviceID:     req.DeviceID,
	})
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewAuthResponse(output.AccessToken, output.RefreshToken, output.User))
}

func (h *UserHandler) Logout(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "invalid user id in context"})
		return
	}

	if err := h.userService.Logout(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrLoginAlreadyExists):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUserInactive):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrSessionNotFound):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
	}
}
