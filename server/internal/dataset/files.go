package dataset

import (
	"context"
	"io"
	storage "robotdatahub/server/internal/platform/storage"
)

type FileService struct {
	datasets *Service
	objects  *storage.ObjectStorage
}

func NewFileService(datasets *Service, objects *storage.ObjectStorage) *FileService {
	return &FileService{datasets: datasets, objects: objects}
}
func (s *FileService) Upload(ctx context.Context, id, name string, r io.Reader, size int64) (string, error) {
	if _, err := s.datasets.Get(ctx, id); err != nil {
		return "", err
	}
	key, err := storage.SafeObjectKey(id, name)
	if err != nil {
		return "", err
	}
	if err = s.objects.Put(ctx, key, r, size, "application/octet-stream"); err != nil {
		return "", err
	}
	return key, nil
}
func (s *FileService) URL(ctx context.Context, id, name string) (string, error) {
	if _, err := s.datasets.Get(ctx, id); err != nil {
		return "", err
	}
	key, err := storage.SafeObjectKey(id, name)
	if err != nil {
		return "", err
	}
	return s.objects.Presign(ctx, key)
}
