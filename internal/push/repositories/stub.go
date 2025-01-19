package repositories

import (
	"github.com/hoitek/Maja-Service/internal/push/domain"
	"github.com/hoitek/Maja-Service/internal/push/models"
)

type PushRepositoryStub struct {
	Pushes []*domain.Push
}

type pushTestCondition struct {
	HasError bool
}

var UserTestCondition *pushTestCondition = &pushTestCondition{}

func NewPushRepositoryStub() *PushRepositoryStub {
	return &PushRepositoryStub{
		Pushes: []*domain.Push{
			{
				ID: 1,
			},
		},
	}
}

func (r *PushRepositoryStub) Query(dataModel *models.PushesQueryRequestParams) ([]*domain.Push, error) {
	var pushes []*domain.Push
	for _, v := range r.Pushes {
		if v.ID == uint(dataModel.ID) {
			pushes = append(pushes, v)
			break
		}
	}
	return pushes, nil
}

func (r *PushRepositoryStub) Count(dataModel *models.PushesQueryRequestParams) (int64, error) {
	var pushes []*domain.Push
	for _, v := range r.Pushes {
		if v.ID == uint(dataModel.ID) {
			pushes = append(pushes, v)
			break
		}
	}
	return int64(len(pushes)), nil
}

func (r *PushRepositoryStub) Migrate() {
	// do stuff
}

func (r *PushRepositoryStub) Seed() {
	// do stuff
}

func (r *PushRepositoryStub) Create(payload *models.PushesCreateRequestBody) (*domain.Push, error) {
	panic("implement me")
}

func (r *PushRepositoryStub) GetPushesByIds(ids []int64) ([]*domain.Push, error) {
	panic("implement me")
}
