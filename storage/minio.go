package storage

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/allegro/bigcache/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"time"
)

// MinIO is the struct for MinIO
type MinIO struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Client          *minio.Client
	Cache           *bigcache.BigCache
}

var MinIOStorage *MinIO

// NewMinIO returns a new MinIO instance
func NewMinIO(endpoint string, accessKeyID string, secretAccessKey string) *MinIO {
	MinIOStorage = &MinIO{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	return MinIOStorage
}

// Connect connects to MinIO
func (m *MinIO) Connect() (*minio.Client, error) {
	// Cache items expire after 30 days
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(24*30*time.Hour))
	if err != nil {
		return nil, err
	}

	client, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKeyID, m.SecretAccessKey, ""),
		Secure: true,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	m.Client = client
	m.Cache = cache

	return client, nil
}

// CreateBucket creates a new bucket
func (m *MinIO) CreateBucket(ctx context.Context, bucketName string) error {
	err := m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

// CheckBucketExists checks if a bucket exists
func (m *MinIO) CheckBucketExists(ctx context.Context, bucketName string) (bool, error) {
	found, err := m.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, err
	}
	return found, nil
}

// UploadFile uploads a file to MinIO
func (m *MinIO) UploadFile(ctx context.Context, bucketName string, objectName string, filePath string) (*minio.UploadInfo, error) {
	info, err := m.Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// UploadFileContent uploads a file content to MinIO
func (m *MinIO) UploadFileContent(ctx context.Context, bucketName string, objectName string, content []byte, userTags map[string]string) (*minio.UploadInfo, error) {
	// Create an io.Reader from the content byte slice.
	reader := bytes.NewReader(content)

	// Upload the object using the io.Reader.
	info, err := m.Client.PutObject(ctx, bucketName, objectName, reader, int64(len(content)), minio.PutObjectOptions{
		UserTags: userTags,
	})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// DownloadFile downloads a file from MinIO
func (m *MinIO) DownloadFile(ctx context.Context, bucketName string, objectName string, filePath string) error {
	err := m.Client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// RemoveFile removes a file from MinIO
func (m *MinIO) RemoveFile(ctx context.Context, bucketName string, objectName string) error {
	err := m.Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// RemoveBucket removes a bucket from MinIO
func (m *MinIO) RemoveBucket(ctx context.Context, bucketName string) error {
	err := m.Client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return err
	}
	return nil
}

// ListBuckets lists all buckets from MinIO
func (m *MinIO) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	buckets, err := m.Client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	return buckets, nil
}

// ListObjects lists all objects from MinIO
func (m *MinIO) ListObjects(ctx context.Context, bucketName string) ([]minio.ObjectInfo, error) {
	objects := make([]minio.ObjectInfo, 0)
	objectCh := m.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

// ListObjectsWithPrefix lists all objects with prefix from MinIO
func (m *MinIO) ListObjectsWithPrefix(ctx context.Context, bucketName string, prefix string) ([]minio.ObjectInfo, error) {
	objects := make([]minio.ObjectInfo, 0)
	objectCh := m.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix: prefix,
	})
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

// GetContext returns a new context
func (m *MinIO) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*1)
}

// StatObject returns object info
func (m *MinIO) StatObject(ctx context.Context, bucketName string, objectName string) (minio.ObjectInfo, error) {
	objectInfo, err := m.Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, err
	}
	return objectInfo, nil
}

// GetObject returns object
func (m *MinIO) GetObject(ctx context.Context, bucketName string, objectName string) (*minio.Object, error) {
	object, err := m.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

// DeleteFilesByTag deletes files by tag
func (m *MinIO) DeleteFilesByTag(ctx context.Context, bucketName string, tagName string, tagValue string) error {
	// Get all objects that have the tag
	objects, err := m.ListObjects(ctx, bucketName)
	if err != nil {
		return err
	}

	// Loop through the objects and delete them
	for _, object := range objects {
		if object.UserTags[tagName] == tagValue {
			err = m.RemoveFile(ctx, bucketName, object.Key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
