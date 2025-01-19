package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/license/domain"
	"github.com/hoitek/Maja-Service/internal/license/models"
)

type LicenseRepositoryStub struct {
	Licenses []*domain.License
}

type licenseTestCondition struct {
	HasError bool
}

var UserTestCondition *licenseTestCondition = &licenseTestCondition{}

func NewLicenseRepositoryStub() *LicenseRepositoryStub {
	return &LicenseRepositoryStub{
		Licenses: []*domain.License{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *LicenseRepositoryStub) Query(dataModel *models.LicensesQueryRequestParams) ([]*domain.License, error) {
	var licenses []*domain.License
	for _, v := range r.Licenses {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			licenses = append(licenses, v)
			break
		}
	}
	return licenses, nil
}

func (r *LicenseRepositoryStub) Count(dataModel *models.LicensesQueryRequestParams) (int64, error) {
	var licenses []*domain.License
	for _, v := range r.Licenses {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			licenses = append(licenses, v)
			break
		}
	}
	return int64(len(licenses)), nil
}

func (r *LicenseRepositoryStub) Migrate() {
	// do stuff
}

func (r *LicenseRepositoryStub) Seed() {
	// do stuff
}

func (r *LicenseRepositoryStub) Create(payload *models.LicensesCreateRequestBody) (*domain.License, error) {
	panic("implement me")
}

func (r *LicenseRepositoryStub) Delete(payload *models.LicensesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *LicenseRepositoryStub) Update(payload *models.LicensesCreateRequestBody, id int64) (*domain.License, error) {
	panic("implement me")
}

func (r *LicenseRepositoryStub) GetLicensesByIds(ids []int64) ([]*domain.License, error) {
	panic("implement me")
}
