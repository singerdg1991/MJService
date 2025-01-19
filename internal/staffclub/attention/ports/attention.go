package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/models"
)

type AttentionService interface {
	Query(dataModel *models.AttentionsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.AttentionsCreateRequestBody) (*domain.Attention, error)
	Delete(payload *models.AttentionsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.AttentionsCreateRequestBody, id int64) (*domain.Attention, error)
	GetAttentionsByIds(ids []int64) ([]*domain.Attention, error)
	FindByID(id int64) (*domain.Attention, error)
	FindByAttentionNumber(attentionNumber int64) (*domain.Attention, error)
}

type AttentionRepositoryPostgresDB interface {
	Query(dataModel *models.AttentionsQueryRequestParams) ([]*domain.Attention, error)
	Count(dataModel *models.AttentionsQueryRequestParams) (int64, error)
	Create(payload *models.AttentionsCreateRequestBody) (*domain.Attention, error)
	Delete(payload *models.AttentionsDeleteRequestBody) ([]int64, error)
	Update(payload *models.AttentionsCreateRequestBody, id int64) (*domain.Attention, error)
	GetAttentionsByIds(ids []int64) ([]*domain.Attention, error)
}
