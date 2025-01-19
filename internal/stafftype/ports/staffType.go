package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/stafftype/domain"
	"github.com/hoitek/Maja-Service/internal/stafftype/models"
)

type StaffTypeService interface {
	Query(dataModel *models.StaffTypesQueryRequestParams) (*restypes.QueryResponse, error)
	GetStaffTypesByIds(ids []int64) ([]*domain.StaffType, error)
}

type StaffTypeRepositoryPostgresDB interface {
	Query(dataModel *models.StaffTypesQueryRequestParams) ([]*domain.StaffType, error)
	Count(dataModel *models.StaffTypesQueryRequestParams) (int64, error)
	GetStaffTypesByIds(ids []int64) ([]*domain.StaffType, error)
}
