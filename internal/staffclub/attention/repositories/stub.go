package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/staffclub/attention/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/models"
)

type AttentionRepositoryStub struct {
	Attentions []*domain.Attention
}

type attentionTestCondition struct {
	HasError bool
}

var UserTestCondition *attentionTestCondition = &attentionTestCondition{}

func NewAttentionRepositoryStub() *AttentionRepositoryStub {
	return &AttentionRepositoryStub{
		Attentions: []*domain.Attention{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *AttentionRepositoryStub) Query(dataModel *models.AttentionsQueryRequestParams) ([]*domain.Attention, error) {
	var attentions []*domain.Attention
	for _, v := range r.Attentions {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			attentions = append(attentions, v)
			break
		}
	}
	return attentions, nil
}

func (r *AttentionRepositoryStub) Count(dataModel *models.AttentionsQueryRequestParams) (int64, error) {
	var attentions []*domain.Attention
	for _, v := range r.Attentions {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			attentions = append(attentions, v)
			break
		}
	}
	return int64(len(attentions)), nil
}

func (r *AttentionRepositoryStub) Migrate() {
	// do stuff
}

func (r *AttentionRepositoryStub) Seed() {
	// do stuff
}

func (r *AttentionRepositoryStub) Create(payload *models.AttentionsCreateRequestBody) (*domain.Attention, error) {
	panic("implement me")
}

func (r *AttentionRepositoryStub) Delete(payload *models.AttentionsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *AttentionRepositoryStub) Update(payload *models.AttentionsCreateRequestBody, id int64) (*domain.Attention, error) {
	panic("implement me")
}

func (r *AttentionRepositoryStub) GetAttentionsByIds(ids []int64) ([]*domain.Attention, error) {
	panic("implement me")
}
