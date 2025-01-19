package repositories

import (
	"github.com/hoitek/Maja-Service/internal/keikkala/domain"
	"github.com/hoitek/Maja-Service/internal/keikkala/models"
)

type KeikkalaRepositoryStub struct {
	Keikkalas []*domain.Keikkala
}

type keikkalaTestCondition struct {
	HasError bool
}

var UserTestCondition *keikkalaTestCondition = &keikkalaTestCondition{}

func NewKeikkalaRepositoryStub() *KeikkalaRepositoryStub {
	return &KeikkalaRepositoryStub{
		Keikkalas: []*domain.Keikkala{
			{
				ID: 1,
			},
		},
	}
}

func (r *KeikkalaRepositoryStub) Query(dataModel *models.KeikkalasQueryRequestParams) ([]*domain.Keikkala, error) {
	var keikkalas []*domain.Keikkala
	for _, v := range r.Keikkalas {
		if v.ID == uint(dataModel.ID) {
			keikkalas = append(keikkalas, v)
			break
		}
	}
	return keikkalas, nil
}

func (r *KeikkalaRepositoryStub) Count(dataModel *models.KeikkalasQueryRequestParams) (int64, error) {
	var keikkalas []*domain.Keikkala
	for _, v := range r.Keikkalas {
		if v.ID == uint(dataModel.ID) {
			keikkalas = append(keikkalas, v)
			break
		}
	}
	return int64(len(keikkalas)), nil
}

func (r *KeikkalaRepositoryStub) Migrate() {
	// do stuff
}

func (r *KeikkalaRepositoryStub) Seed() {
	// do stuff
}

func (r *KeikkalaRepositoryStub) Create(payload *models.KeikkalasCreateRequestBody) (*domain.Keikkala, error) {
	panic("implement me")
}

func (r *KeikkalaRepositoryStub) Delete(payload *models.KeikkalasDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *KeikkalaRepositoryStub) GetKeikkalaShiftsByIds(ids []int64) ([]*domain.Keikkala, error) {
	panic("implement me")
}

func (r *KeikkalaRepositoryStub) QueryShiftStatistics(queries *models.KeikkalasQueryShiftStatisticsRequestParams) (int, int, int, error) {
	panic("implement me")
}
