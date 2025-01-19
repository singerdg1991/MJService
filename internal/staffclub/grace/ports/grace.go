package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace/models"
)

type GraceService interface {
	Query(dataModel *models.GracesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.GracesCreateRequestBody) (*domain.Grace, error)
	Delete(payload *models.GracesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.GracesCreateRequestBody, id int64) (*domain.Grace, error)
	GetGracesByIds(ids []int64) ([]*domain.Grace, error)
	FindByID(id int64) (*domain.Grace, error)
	FindByGraceNumber(graceNumber int64) (*domain.Grace, error)
}

type GraceRepositoryPostgresDB interface {
	Query(dataModel *models.GracesQueryRequestParams) ([]*domain.Grace, error)
	Count(dataModel *models.GracesQueryRequestParams) (int64, error)
	Create(payload *models.GracesCreateRequestBody) (*domain.Grace, error)
	Delete(payload *models.GracesDeleteRequestBody) ([]int64, error)
	Update(payload *models.GracesCreateRequestBody, id int64) (*domain.Grace, error)
	GetGracesByIds(ids []int64) ([]*domain.Grace, error)
}
