package ports

import (
	"github.com/hoitek/Govalidity/govalidityt"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
)

type S3Service interface {
	GetCacheKey(bucketName string, fileName string) string
	Cache(bucketName string, fileName string, entry []byte) error
	UploadFile(bucketName string, file *govalidityt.File, id int64) (*types.UploadMetadata, error)
	UploadFiles(bucketName string, files []*govalidityt.File, id int64) ([]*types.UploadMetadata, []error)
	DeleteFilesByTag(bucketName string, tagName string, tagValue string) error
	DeleteFilesByTagID(bucketName string, id int64) error
}
