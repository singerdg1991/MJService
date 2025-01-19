package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/s3/ports"
	"github.com/hoitek/Maja-Service/storage"
	"io"
	"log"
	"path/filepath"
	"sync"
	"time"
)

type S3Service struct {
	MinIOStorage *storage.MinIO
}

func NewS3Service(s *storage.MinIO) ports.S3Service {
	return &S3Service{
		MinIOStorage: s,
	}
}

func (s *S3Service) MakeUniqueUploadFileName(bucketName string, fileName string) string {
	// Get the file name with extension
	baseName := filepath.Base(fileName)

	// Get the file extension with the dot
	ext := filepath.Ext(baseName)

	// Get the file name without extension
	fileNameWithoutExt := baseName[:len(baseName)-len(ext)]

	// Now we need to make sure that the file name is unique
	newFileName := fmt.Sprintf("%s-%s-%d", bucketName, fileNameWithoutExt, time.Now().UnixNano())

	// Create a new SHA-256 hash object
	hasher := md5.New()
	_, err := hasher.Write([]byte(newFileName))
	if err != nil {
		return fmt.Sprintf("%s%s", newFileName, ext)
	}

	// Calculate the hash sum
	hashSum := hasher.Sum(nil)

	// Convert the hash sum to a hexadecimal string
	hashHex := hex.EncodeToString(hashSum[:8])
	hashHex = fmt.Sprintf("%s%d", hashHex, time.Now().UnixNano())

	return fmt.Sprintf("%s%s", hashHex, ext)
}

func (s *S3Service) GetCacheKey(bucketName string, fileName string) string {
	return fmt.Sprintf("%s.%s", bucketName, fileName)
}

func (s *S3Service) Cache(bucketName string, fileName string, entry []byte) error {
	cacheKey := s.GetCacheKey(bucketName, fileName)
	return s.MinIOStorage.Cache.Set(cacheKey, entry)
}

func (s *S3Service) UploadFile(bucketName string, file *govalidityt.File, id int64) (*types.UploadMetadata, error) {
	fileName := s.MakeUniqueUploadFileName(bucketName, file.Header.Filename)
	content, err := io.ReadAll(*file.File)
	if err != nil {
		return nil, err
	}

	// Upload the file with FPutObject
	info, err := s.MinIOStorage.UploadFileContent(context.Background(), bucketName, fileName, content, map[string]string{
		"ID": fmt.Sprintf("%d", id),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Uploaded file %s successfully to bucket %s\n", fileName, bucketName)

	// Cache the file
	err = s.Cache(bucketName, fileName, content)
	if err != nil {
		log.Printf("Failed to cache file %s to bucket %s\n", fileName, bucketName)
	}

	return &types.UploadMetadata{
		FileName: fileName,
		FileSize: info.Size,
	}, nil
}

func (s *S3Service) UploadFiles(bucketName string, files []*govalidityt.File, id int64) ([]*types.UploadMetadata, []error) {
	type result struct {
		metadata *types.UploadMetadata
		err      error
	}
	var (
		errors     []error
		wg         sync.WaitGroup
		metadata   []*types.UploadMetadata
		resultChan = make(chan *result, len(files))
	)
	wg.Add(len(files))
	for _, attachment := range files {
		go func(f *govalidityt.File) {
			defer wg.Done()
			m, err := s.UploadFile(bucketName, f, id)
			if err == nil {
				resultChan <- &result{
					metadata: m,
					err:      nil,
				}
			} else {
				resultChan <- &result{
					metadata: nil,
					err:      err,
				}
			}
		}(attachment)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.err != nil {
			errors = append(errors, result.err)
		} else {
			metadata = append(metadata, result.metadata)
		}
	}

	if len(errors) > 0 {
		log.Printf("UploadFiles errors: %v in bucket %s", errors, bucketName)
	}
	return metadata, errors
}

func (s *S3Service) DeleteFilesByTag(bucketName string, tagName string, tagValue string) error {
	return s.MinIOStorage.DeleteFilesByTag(context.Background(), bucketName, tagName, tagValue)
}

func (s *S3Service) DeleteFilesByTagID(bucketName string, id int64) error {
	return s.MinIOStorage.DeleteFilesByTag(context.Background(), bucketName, "ID", fmt.Sprintf("%d", id))
}
