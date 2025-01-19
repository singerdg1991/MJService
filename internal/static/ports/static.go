package ports

import "github.com/minio/minio-go/v7"

type StaticService interface {
	CheckBucketExists(bucketName string) (bool, error)
	StatObject(bucketName string, objectName string) (minio.ObjectInfo, error)
	GetObject(bucketName string, objectName string) (*minio.Object, error)
	SetCache(key string, value []byte) error
	GetCache(key string) ([]byte, error)
}
