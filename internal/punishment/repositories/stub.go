package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/punishment/domain"
	"github.com/hoitek/Maja-Service/internal/punishment/models"
)

type PunishmentRepositoryStub struct {
	Punishments []*domain.Punishment
}

type punishmentTestCondition struct {
	HasError bool
}

var UserTestCondition *punishmentTestCondition = &punishmentTestCondition{}

func NewPunishmentRepositoryStub() *PunishmentRepositoryStub {
	return &PunishmentRepositoryStub{
		Punishments: []*domain.Punishment{
			{
				ID:   1,
				Name: "test",
			},
		},
	}
}

func (r *PunishmentRepositoryStub) Query(dataModel *models.PunishmentsQueryRequestParams) ([]*domain.Punishment, error) {
	var punishments []*domain.Punishment
	for _, v := range r.Punishments {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			punishments = append(punishments, v)
			break
		}
	}
	return punishments, nil
}

func (r *PunishmentRepositoryStub) Count(dataModel *models.PunishmentsQueryRequestParams) (int64, error) {
	var punishments []*domain.Punishment
	for _, v := range r.Punishments {
		if v.ID == uint(dataModel.ID) ||
			v.Name == fmt.Sprintf("%v", dataModel.Filters.Name) {
			punishments = append(punishments, v)
			break
		}
	}
	return int64(len(punishments)), nil
}

func (r *PunishmentRepositoryStub) Migrate() {
	// do stuff
}

func (r *PunishmentRepositoryStub) Seed() {
	// do stuff
}

func (r *PunishmentRepositoryStub) Create(payload *models.PunishmentsCreateRequestBody) (*domain.Punishment, error) {
	panic("implement me")
}

func (r *PunishmentRepositoryStub) Delete(payload *models.PunishmentsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *PunishmentRepositoryStub) Update(payload *models.PunishmentsCreateRequestBody, id int64) (*domain.Punishment, error) {
	panic("implement me")
}

func (r *PunishmentRepositoryStub) GetPunishmentsByIds(ids []int64) ([]*domain.Punishment, error) {
	panic("implement me")
}
