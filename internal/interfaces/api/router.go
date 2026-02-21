package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prodonik/bank_app/internal/infrastructure/auth"
	"github.com/prodonik/bank_app/internal/interfaces/api/handler"
	"github.com/prodonik/bank_app/internal/interfaces/api/middleware"
)

func NewRouter(userHandler *handler.UserHandler, jwtService *auth.JWTService) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.POST("/refresh", userHandler.Refresh)

			authenticated := users.Group("")
			authenticated.Use(middleware.AuthMiddleware(jwtService))
			{
				authenticated.POST("/logout", userHandler.Logout)
			}
		}
	}

	return r
}
