package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/archive/domain"
	"github.com/hoitek/Maja-Service/internal/archive/models"
)

type ArchiveService interface {
	Query(dataModel *models.ArchivesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.ArchivesCreateRequestBody) (*domain.Archive, error)
	Delete(payload *models.ArchivesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.ArchivesCreateRequestBody, id int64) (*domain.Archive, error)
	GetArchivesByIds(ids []int64) ([]*domain.Archive, error)
	FindByID(id int64) (*domain.Archive, error)
	FindByTitle(title string) (*domain.Archive, error)
	UpdateAttachments(attachments []*types.UploadMetadata, id int64) (*domain.Archive, error)
}

type ArchiveRepositoryPostgresDB interface {
	Query(dataModel *models.ArchivesQueryRequestParams) ([]*domain.Archive, error)
	Count(dataModel *models.ArchivesQueryRequestParams) (int64, error)
	Create(payload *models.ArchivesCreateRequestBody) (*domain.Archive, error)
	Delete(payload *models.ArchivesDeleteRequestBody) ([]int64, error)
	Update(payload *models.ArchivesCreateRequestBody, id int64) (*domain.Archive, error)
	GetArchivesByIds(ids []int64) ([]*domain.Archive, error)
	UpdateAttachments(attachments []*types.UploadMetadata, id int64) (*domain.Archive, error)
}
