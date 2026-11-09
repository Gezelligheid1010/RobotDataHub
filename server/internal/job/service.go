package job

import (
	"context"
	"encoding/json"
	"robotdatahub/server/internal/dataset"
	"time"
)

type Service struct {
	repo     *Repository
	datasets *dataset.Service
}

func NewService(repo *Repository, datasets *dataset.Service) *Service {
	return &Service{repo: repo, datasets: datasets}
}
func (s *Service) Create(ctx context.Context, datasetID string) (Job, error) {
	if _, err := s.datasets.Get(ctx, datasetID); err != nil {
		return Job{}, err
	}
	return s.repo.CreateJob(ctx, datasetID)
}
func (s *Service) List(ctx context.Context) ([]Job, error) { return s.repo.ListJobs(ctx) }
func (s *Service) Claim(ctx context.Context) (Job, error)  { return s.repo.ClaimJob(ctx) }
func (s *Service) Complete(ctx context.Context, id string, result json.RawMessage, jobError string) error {
	return s.repo.CompleteJob(ctx, id, result, jobError)
}
func (s *Service) Recover(ctx context.Context, age time.Duration) error {
	return s.repo.RecoverJobs(ctx, age)
}
