package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/prodonik/bank_app/internal/domain/user"
	"github.com/prodonik/bank_app/internal/infrastructure/database/sqlc"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		queries: sqlc.New(db),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	row, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		FullName:     u.FullName,
		Role:         u.Role,
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		Status:       u.Status,
	})
	if err != nil {
		return nil, err
	}
	return mapSqlcUser(row), nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	row, err := r.queries.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return mapSqlcUser(row), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return mapSqlcUser(row), nil
}

func mapSqlcUser(u sqlc.User) *user.User {
	return &user.User{
		ID:           u.ID,
		FullName:     u.FullName,
		Role:         u.Role,
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		Status:       u.Status,
		CreatedAt:    u.CreatedAt,
	}
}
