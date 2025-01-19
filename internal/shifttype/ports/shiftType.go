package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/shifttype/domain"
	"github.com/hoitek/Maja-Service/internal/shifttype/models"
)

type ShiftTypeService interface {
	Query(dataModel *models.ShiftTypesQueryRequestParams) (*restypes.QueryResponse, error)
	GetShiftTypesByIds(ids []int64) ([]*domain.ShiftType, error)
}

type ShiftTypeRepositoryPostgresDB interface {
	Query(dataModel *models.ShiftTypesQueryRequestParams) ([]*domain.ShiftType, error)
	Count(dataModel *models.ShiftTypesQueryRequestParams) (int64, error)
	GetShiftTypesByIds(ids []int64) ([]*domain.ShiftType, error)
}
