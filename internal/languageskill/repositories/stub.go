package repositories

import (
	"fmt"
	"github.com/hoitek/Maja-Service/internal/languageskill/domain"
	"github.com/hoitek/Maja-Service/internal/languageskill/models"
)

type LanguageSkillRepositoryStub struct {
	LanguageSkills []*domain.LanguageSkill
}

type languageskillTestCondition struct {
	HasError bool
}

var UserTestCondition *languageskillTestCondition = &languageskillTestCondition{}

func NewLanguageSkillRepositoryStub() *LanguageSkillRepositoryStub {
	return &LanguageSkillRepositoryStub{
		LanguageSkills: []*domain.LanguageSkill{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *LanguageSkillRepositoryStub) Query(dataModel *models.LanguageSkillsQueryRequestParams) ([]*domain.LanguageSkill, error) {
	var languageskills []*domain.LanguageSkill
	for _, v := range r.LanguageSkills {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			languageskills = append(languageskills, v)
			break
		}
	}
	return languageskills, nil
}

func (r *LanguageSkillRepositoryStub) Count(dataModel *models.LanguageSkillsQueryRequestParams) (int64, error) {
	var languageskills []*domain.LanguageSkill
	for _, v := range r.LanguageSkills {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			languageskills = append(languageskills, v)
			break
		}
	}
	return int64(len(languageskills)), nil
}

func (r *LanguageSkillRepositoryStub) Migrate() {
	// do stuff
}

func (r *LanguageSkillRepositoryStub) Seed() {
	// do stuff
}

func (r *LanguageSkillRepositoryStub) Create(payload *models.LanguageSkillsCreateRequestBody) (*domain.LanguageSkill, error) {
	panic("implement me")
}

func (r *LanguageSkillRepositoryStub) Delete(payload *models.LanguageSkillsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *LanguageSkillRepositoryStub) Update(payload *models.LanguageSkillsCreateRequestBody, id int64) (*domain.LanguageSkill, error) {
	panic("implement me")
}

func (r *LanguageSkillRepositoryStub) GetLanguageSkillsByIds(ids []int64) ([]*domain.LanguageSkill, error) {
	panic("implement me")
}
