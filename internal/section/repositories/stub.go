package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/section/domain"
	"github.com/hoitek/Maja-Service/internal/section/models"
)

type SectionRepositoryStub struct {
	Sections []*domain.Section
}

type sectionTestCondition struct {
	HasError bool
}

var UserTestCondition *sectionTestCondition = &sectionTestCondition{}

func NewSectionRepositoryStub() *SectionRepositoryStub {
	return &SectionRepositoryStub{
		Sections: []*domain.Section{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *SectionRepositoryStub) Query(dataModel *models.SectionsQueryRequestParams) ([]*domain.Section, error) {
	var sections []*domain.Section
	for _, v := range r.Sections {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			sections = append(sections, v)
			break
		}
	}
	return sections, nil
}

func (r *SectionRepositoryStub) Count(dataModel *models.SectionsQueryRequestParams) (int64, error) {
	var sections []*domain.Section
	for _, v := range r.Sections {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			sections = append(sections, v)
			break
		}
	}
	return int64(len(sections)), nil
}

func (r *SectionRepositoryStub) Migrate() {
	// do stuff
}

func (r *SectionRepositoryStub) Seed(sections []*domain.Section) error {
	// do stuff
	return nil
}

func (r *SectionRepositoryStub) Create(payload *models.SectionsCreateRequestBody) (*domain.Section, error) {
	panic("implement me")
}

func (r *SectionRepositoryStub) Delete(payload *models.SectionsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *SectionRepositoryStub) Update(payload *models.SectionsCreateRequestBody, id int) (*domain.Section, error) {
	panic("implement me")
}

func (r *SectionRepositoryStub) GetSectionsByIds(ids []int64) ([]*domain.Section, error) {
	panic("implement me")
}
