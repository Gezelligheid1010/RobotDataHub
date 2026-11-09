package episode

import "context"

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{repo: r} }
func (s *Service) List(ctx context.Context, id string) ([]Episode, error) {
	return s.repo.ListEpisodes(ctx, id)
}
func (s *Service) Import(ctx context.Context, id string, raw []byte) (int, int, error) {
	eps, total, fps, dim, err := ParseManifest(raw)
	if err != nil {
		return 0, 0, err
	}
	if err = s.repo.ImportEpisodes(ctx, id, eps, total, fps, dim); err != nil {
		return 0, 0, err
	}
	return len(eps), total, nil
}
