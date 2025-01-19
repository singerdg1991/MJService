package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/ability/domain"
	"github.com/hoitek/Maja-Service/internal/ability/models"
)

type AbilityService interface {
	Query(dataModel *models.AbilitiesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.AbilitiesCreateRequestBody) (*domain.Ability, error)
	Delete(payload *models.AbilitiesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.AbilitiesCreateRequestBody, name string) (*domain.Ability, error)
	GetAbilitiesByIds(ids []int64) ([]*domain.Ability, error)
}

type AbilityRepositoryPostgresDB interface {
	Query(dataModel *models.AbilitiesQueryRequestParams) ([]*domain.Ability, error)
	Count(dataModel *models.AbilitiesQueryRequestParams) (int64, error)
	Create(payload *models.AbilitiesCreateRequestBody) (*domain.Ability, error)
	Delete(payload *models.AbilitiesDeleteRequestBody) ([]int64, error)
	Update(payload *models.AbilitiesCreateRequestBody, name string) (*domain.Ability, error)
	GetAbilitiesByIds(ids []int64) ([]*domain.Ability, error)
}
