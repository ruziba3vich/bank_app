package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prodonik/bank_app/internal/infrastructure/auth"
	"github.com/prodonik/bank_app/internal/interfaces/api/handler"
	"github.com/prodonik/bank_app/internal/interfaces/api/middleware"
)

func NewRouter(userHandler *handler.UserHandler, jwtService *auth.JWTService) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
				authenticated.GET("/me", userHandler.GetMe)
				authenticated.GET("", userHandler.GetAll)
				authenticated.GET("/:id", userHandler.GetByID)
				authenticated.PUT("/:id", userHandler.Update)
				authenticated.DELETE("/:id", userHandler.Delete)
			}
		}
	}

	return r
}
