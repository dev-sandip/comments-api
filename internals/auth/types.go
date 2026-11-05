package auth

import (
	"context"

	repo "github.com/dev-sandip/comments-api/internals/adapters/postgres/sqlc"
)

type Service interface {
	GetUserByEmail(ctx context.Context, email string) (repo.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) (repo.User, error)
}

type svc struct {
	repo repo.Querier
}

type handler struct {
	service Service
}

type CreateUserParams struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginParams struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"password123"`
}
