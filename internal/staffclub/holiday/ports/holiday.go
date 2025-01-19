package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/domain"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/models"
)

type HolidayService interface {
	Query(dataModel *models.HolidaysQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.HolidaysCreateRequestBody) (*domain.Holiday, error)
	Delete(payload *models.HolidaysDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.HolidaysCreateRequestBody, id int64) (*domain.Holiday, error)
	GetHolidaysByIds(ids []int64) ([]*domain.Holiday, error)
	FindByID(id int64) (*domain.Holiday, error)
	UpdateStatus(payload *models.HolidaysUpdateStatusRequestBody, id int64) (*domain.Holiday, error)
}

type HolidayRepositoryPostgresDB interface {
	Query(dataModel *models.HolidaysQueryRequestParams) ([]*domain.Holiday, error)
	Count(dataModel *models.HolidaysQueryRequestParams) (int64, error)
	Create(payload *models.HolidaysCreateRequestBody) (*domain.Holiday, error)
	Delete(payload *models.HolidaysDeleteRequestBody) ([]int64, error)
	Update(payload *models.HolidaysCreateRequestBody, id int64) (*domain.Holiday, error)
	GetHolidaysByIds(ids []int64) ([]*domain.Holiday, error)
	UpdateStatus(payload *models.HolidaysUpdateStatusRequestBody, id int64) (*domain.Holiday, error)
}
