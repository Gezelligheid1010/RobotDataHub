package storage

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/url"
	"path"
	"strings"
	"time"
)

type ObjectStorage struct {
	Client *minio.Client
	Bucket string
}

func NewObjectStorage(endpoint, access, secret, bucket string, ssl bool) (*ObjectStorage, error) {
	c, err := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(access, secret, ""), Secure: ssl})
	if err != nil {
		return nil, err
	}
	return &ObjectStorage{Client: c, Bucket: bucket}, nil
}
func (o *ObjectStorage) EnsureBucket(ctx context.Context) error {
	exists, err := o.Client.BucketExists(ctx, o.Bucket)
	if err != nil {
		return err
	}
	if !exists {
		return o.Client.MakeBucket(ctx, o.Bucket, minio.MakeBucketOptions{})
	}
	return nil
}
func (o *ObjectStorage) Put(ctx context.Context, key string, src io.Reader, size int64, contentType string) error {
	_, err := o.Client.PutObject(ctx, o.Bucket, key, src, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}
func (o *ObjectStorage) Read(ctx context.Context, key string) ([]byte, error) {
	obj, err := o.Client.GetObject(ctx, o.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(io.LimitReader(obj, 16<<20))
}
func (o *ObjectStorage) Presign(ctx context.Context, key string) (string, error) {
	u, err := o.Client.PresignedGetObject(ctx, o.Bucket, key, 15*time.Minute, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func SafeObjectKey(datasetID, filename string) (string, error) {
	filename = strings.ReplaceAll(filename, "\\", "/")
	filename = path.Clean(filename)
	if filename == "." || strings.HasPrefix(filename, "/") || filename == ".." || strings.HasPrefix(filename, "../") || strings.Contains(filename, ":") {
		return "", errors.New("unsafe file path")
	}
	return "datasets/" + datasetID + "/" + filename, nil
}
