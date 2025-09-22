package health

import (
	"context"
	"fmt"
	"prodcrud/internal/repository/health"
)

type Service struct {
	repo *health.Repo
}

func NewService(repo *health.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Check(ctx context.Context) error {
	if err := s.repo.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	return nil
}
