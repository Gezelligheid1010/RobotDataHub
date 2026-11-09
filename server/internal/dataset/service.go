package dataset

import (
	"context"
	"errors"
	"strings"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, req CreateRequest) (Dataset, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return Dataset{}, errors.New("name required")
	}
	return s.repo.CreateDataset(ctx, req.Name, req.Description, req.Source)
}
func (s *Service) List(ctx context.Context) ([]Dataset, error) { return s.repo.ListDatasets(ctx) }
func (s *Service) Get(ctx context.Context, id string) (Dataset, error) {
	return s.repo.GetDataset(ctx, id)
}
func (s *Service) Delete(ctx context.Context, id string) error { return s.repo.DeleteDataset(ctx, id) }
