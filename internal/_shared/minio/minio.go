package minio

import (
	"context"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/minio/minio-go/v7"
	"log"
	"time"
)

// SetupMinIOStorage creates the bucket if it does not exist
func SetupMinIOStorage(bucketName string, m *storage.MinIO) {
	// Check environment if testing or empty, skip
	if config.AppConfig.Environment == "test" || config.AppConfig.Environment == "" {
		return
	}

	// Check if bucket exists and create if not
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	isBucketExists, err := m.CheckBucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Unable to check if bucket exists for %s: %v", bucketName, err)
	}
	if !isBucketExists {
		// Create bucket
		err = m.CreateBucket(ctx, bucketName)
		if err != nil {
			if minioErr, ok := err.(minio.ErrorResponse); ok {
				if minioErr.Code == "BucketAlreadyOwnedByYou" {
					return
				}
			}
			log.Fatalf("Unable to create bucket for %s: %v", bucketName, err)
		}
	}
}
