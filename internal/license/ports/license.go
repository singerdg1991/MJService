package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/license/domain"
	"github.com/hoitek/Maja-Service/internal/license/models"
)

type LicenseService interface {
	Query(dataModel *models.LicensesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.LicensesCreateRequestBody) (*domain.License, error)
	Delete(payload *models.LicensesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.LicensesCreateRequestBody, id int64) (*domain.License, error)
	GetLicensesByIds(ids []int64) ([]*domain.License, error)
	FindByID(id int64) (*domain.License, error)
	FindByName(name string) (*domain.License, error)
}

type LicenseRepositoryPostgresDB interface {
	Query(dataModel *models.LicensesQueryRequestParams) ([]*domain.License, error)
	Count(dataModel *models.LicensesQueryRequestParams) (int64, error)
	Create(payload *models.LicensesCreateRequestBody) (*domain.License, error)
	Delete(payload *models.LicensesDeleteRequestBody) ([]int64, error)
	Update(payload *models.LicensesCreateRequestBody, id int64) (*domain.License, error)
	GetLicensesByIds(ids []int64) ([]*domain.License, error)
}
