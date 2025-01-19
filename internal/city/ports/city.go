package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/city/domain"
	"github.com/hoitek/Maja-Service/internal/city/models"
)

type CityService interface {
	Query(dataModel *models.CitiesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.CitiesCreateRequestBody) (*domain.City, error)
	Delete(payload *models.CitiesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.CitiesCreateRequestBody, name string) (*domain.City, error)
	GetCityByID(id int64) (*domain.City, error)
}

type CityRepositoryPostgresDB interface {
	Query(dataModel *models.CitiesQueryRequestParams) ([]*domain.City, error)
	Count(dataModel *models.CitiesQueryRequestParams) (int64, error)
	Create(payload *models.CitiesCreateRequestBody) (*domain.City, error)
	Delete(payload *models.CitiesDeleteRequestBody) ([]int64, error)
	Update(payload *models.CitiesCreateRequestBody, name string) (*domain.City, error)
}
