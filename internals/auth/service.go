package auth

import (
	"context"
	"database/sql"
	"errors"

	repo "github.com/dev-sandip/comments-api/internals/adapters/postgres/sqlc"
)

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.User{}, ErrUserNotFound
		}
		return repo.User{}, err
	}

	return user, nil
}

func (s *svc) CreateUser(ctx context.Context, params CreateUserParams) (repo.User, error) {
	return s.repo.CreateUser(ctx, repo.CreateUserParams{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	})
}
