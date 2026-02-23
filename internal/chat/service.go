package chat

import "context"

type Service interface {
	Create(ctx context.Context, input CreateMessageInput) (Message, error)
	List(ctx context.Context) ([]Message, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateMessageInput) (Message, error) {
	return s.repo.Create(ctx, input)
}

func (s *service) List(ctx context.Context) ([]Message, error) {
	return s.repo.List(ctx)
}
