package user

import "context"

type Service interface {
	Create(ctx context.Context, input CreateUserInput) (User, error)
	List(ctx context.Context) ([]User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateUserInput) (User, error) {
	return s.repo.Create(ctx, input)
}

func (s *service) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}
