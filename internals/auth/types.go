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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
