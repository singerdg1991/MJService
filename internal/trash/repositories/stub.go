package repositories

import (
	"github.com/hoitek/Maja-Service/internal/trash/domain"
	"github.com/hoitek/Maja-Service/internal/trash/models"
)

type TrashRepositoryStub struct {
	Trashes []*domain.Trash
}

type trashTestCondition struct {
	HasError bool
}

var UserTestCondition *trashTestCondition = &trashTestCondition{}

func NewTrashRepositoryStub() *TrashRepositoryStub {
	return &TrashRepositoryStub{
		Trashes: []*domain.Trash{
			{
				ID: 1,
			},
		},
	}
}

func (r *TrashRepositoryStub) Query(dataModel *models.TrashesQueryRequestParams) ([]*domain.Trash, error) {
	var trashes []*domain.Trash
	panic("implement me")
	return trashes, nil
}

func (r *TrashRepositoryStub) Count(dataModel *models.TrashesQueryRequestParams) (int64, error) {
	var trashes []*domain.Trash
	panic("implement me")
	return int64(len(trashes)), nil
}

func (r *TrashRepositoryStub) Create(payload *models.TrashesCreateRequestBody) (*domain.Trash, error) {
	panic("implement me")
}

func (r *TrashRepositoryStub) Delete(payload *models.TrashesDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}
