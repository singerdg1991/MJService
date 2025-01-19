package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/company/domain"
	"github.com/hoitek/Maja-Service/internal/company/models"
)

type CompanyRepositoryStub struct {
	Companies []*domain.Company
}

type companyTestCondition struct {
	HasError bool
}

var UserTestCondition *companyTestCondition = &companyTestCondition{}

func NewCompanyRepositoryStub() *CompanyRepositoryStub {
	return &CompanyRepositoryStub{
		Companies: []*domain.Company{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *CompanyRepositoryStub) Query(dataModel *models.CompaniesQueryRequestParams) ([]*domain.Company, error) {
	var companies []*domain.Company
	for _, v := range r.Companies {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			companies = append(companies, v)
			break
		}
	}
	return companies, nil
}

func (r *CompanyRepositoryStub) Count(dataModel *models.CompaniesQueryRequestParams) (int64, error) {
	var companies []*domain.Company
	for _, v := range r.Companies {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			companies = append(companies, v)
			break
		}
	}
	return int64(len(companies)), nil
}

func (r *CompanyRepositoryStub) Migrate() {
	// do stuff
}

func (r *CompanyRepositoryStub) Seed() {
	// do stuff
}

func (r *CompanyRepositoryStub) Create(payload *models.CompaniesCreateRequestBody) (*domain.Company, error) {
	panic("implement me")
}

func (r *CompanyRepositoryStub) Delete(payload *models.CompaniesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *CompanyRepositoryStub) Update(payload *models.CompaniesCreateRequestBody, name string) (*domain.Company, error) {
	panic("implement me")
}
