package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/prodonik/bank_app/config"
	_ "github.com/prodonik/bank_app/docs"
	appuser "github.com/prodonik/bank_app/internal/application/user"
	"github.com/prodonik/bank_app/internal/infrastructure/auth"
	"github.com/prodonik/bank_app/internal/infrastructure/database"
	"github.com/prodonik/bank_app/internal/infrastructure/repository"
	"github.com/prodonik/bank_app/internal/interfaces/api"
	"github.com/prodonik/bank_app/internal/interfaces/api/handler"
)

// @title Bank App API
// @version 1.0
// @description Authentication service for the Bank App
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}"
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Run migrations
	m, err := migrate.New("file://internal/infrastructure/database/migrations", cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Database connection
	db, err := database.NewConnection(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Infrastructure
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Application
	userService := appuser.NewService(userRepo, sessionRepo, jwtService)

	// Interfaces
	userHandler := handler.NewUserHandler(userService)
	router := api.NewRouter(userHandler, jwtService)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
