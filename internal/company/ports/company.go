package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/company/domain"
	"github.com/hoitek/Maja-Service/internal/company/models"
)

type CompanyService interface {
	Query(dataModel *models.CompaniesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.CompaniesCreateRequestBody) (*domain.Company, error)
	Delete(payload *models.CompaniesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.CompaniesCreateRequestBody, name string) (*domain.Company, error)
}

type CompanyRepositoryPostgresDB interface {
	Query(dataModel *models.CompaniesQueryRequestParams) ([]*domain.Company, error)
	Count(dataModel *models.CompaniesQueryRequestParams) (int64, error)
	Create(payload *models.CompaniesCreateRequestBody) (*domain.Company, error)
	Delete(payload *models.CompaniesDeleteRequestBody) ([]int64, error)
	Update(payload *models.CompaniesCreateRequestBody, name string) (*domain.Company, error)
}
