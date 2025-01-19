package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/section/domain"
	"github.com/hoitek/Maja-Service/internal/section/models"
)

type SectionService interface {
	Query(dataModel *models.SectionsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.SectionsCreateRequestBody) (*domain.Section, error)
	Delete(payload *models.SectionsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.SectionsCreateRequestBody, id int) (*domain.Section, error)
	GetSectionsByIds(ids []int64) ([]*domain.Section, error)
}

type SectionRepositoryPostgresDB interface {
	Query(dataModel *models.SectionsQueryRequestParams) ([]*domain.Section, error)
	Count(dataModel *models.SectionsQueryRequestParams) (int64, error)
	Create(payload *models.SectionsCreateRequestBody) (*domain.Section, error)
	Delete(payload *models.SectionsDeleteRequestBody) ([]int64, error)
	Update(payload *models.SectionsCreateRequestBody, id int) (*domain.Section, error)
	GetSectionsByIds(ids []int64) ([]*domain.Section, error)
}

type SectionRepositoryMongoDB interface {
	Query(queries *models.SectionsQueryRequestParams) ([]*domain.Section, error)
	Count(queries *models.SectionsQueryRequestParams) (int64, error)
	Create(postgresID int, payload interface{}) (interface{}, error)
	CreateOrUpdate(postgresID int, payload interface{}) (interface{}, error)
	AppendChildren(children interface{}, id int) error
	Update(payload interface{}, id int) error
	Delete(ids []uint) error
}
