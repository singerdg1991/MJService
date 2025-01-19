package service

import (
	"context"
	"github.com/hoitek/Maja-Service/internal/static/ports"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/minio/minio-go/v7"
)

type StaticService struct {
	MinIOStorage *storage.MinIO
}

func NewStaticService(m *storage.MinIO) ports.StaticService {
	return &StaticService{
		MinIOStorage: m,
	}
}

func (s *StaticService) SetCache(key string, value []byte) error {
	return s.MinIOStorage.Cache.Set(key, value)
}

func (s *StaticService) GetCache(key string) ([]byte, error) {
	return s.MinIOStorage.Cache.Get(key)
}

func (s *StaticService) CheckBucketExists(bucketName string) (bool, error) {
	ctx := context.Background()
	return s.MinIOStorage.CheckBucketExists(ctx, bucketName)
}

func (s *StaticService) StatObject(bucketName string, objectName string) (minio.ObjectInfo, error) {
	ctx := context.Background()
	return s.MinIOStorage.StatObject(ctx, bucketName, objectName)
}

func (s *StaticService) GetObject(bucketName string, objectName string) (*minio.Object, error) {
	ctx := context.Background()
	return s.MinIOStorage.GetObject(ctx, bucketName, objectName)
}
