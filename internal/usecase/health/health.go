package health

import (
	"context"
	"prodcrud/internal/repository/health"
)

type Service struct {
	repo *health.Repo
}

func NewService(repo *health.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Check(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
