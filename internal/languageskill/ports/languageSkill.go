package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/languageskill/domain"
	"github.com/hoitek/Maja-Service/internal/languageskill/models"
)

type LanguageSkillService interface {
	Query(dataModel *models.LanguageSkillsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.LanguageSkillsCreateRequestBody) (*domain.LanguageSkill, error)
	Delete(payload *models.LanguageSkillsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.LanguageSkillsCreateRequestBody, id int64) (*domain.LanguageSkill, error)
	GetLanguageSkillsByIds(ids []int64) ([]*domain.LanguageSkill, error)
	FindByID(id int64) (*domain.LanguageSkill, error)
	FindByName(name string) (*domain.LanguageSkill, error)
}

type LanguageSkillRepositoryPostgresDB interface {
	Query(dataModel *models.LanguageSkillsQueryRequestParams) ([]*domain.LanguageSkill, error)
	Count(dataModel *models.LanguageSkillsQueryRequestParams) (int64, error)
	Create(payload *models.LanguageSkillsCreateRequestBody) (*domain.LanguageSkill, error)
	Delete(payload *models.LanguageSkillsDeleteRequestBody) ([]int64, error)
	Update(payload *models.LanguageSkillsCreateRequestBody, id int64) (*domain.LanguageSkill, error)
	GetLanguageSkillsByIds(ids []int64) ([]*domain.LanguageSkill, error)
}
