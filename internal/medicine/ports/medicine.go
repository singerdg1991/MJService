package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/medicine/domain"
	"github.com/hoitek/Maja-Service/internal/medicine/models"
)

type MedicineService interface {
	Query(dataModel *models.MedicinesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.MedicinesCreateRequestBody) (*domain.Medicine, error)
	Delete(payload *models.MedicinesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.MedicinesCreateRequestBody, id int64) (*domain.Medicine, error)
	GetMedicinesByIds(ids []int64) ([]*domain.Medicine, error)
	FindByID(id int64) (*domain.Medicine, error)
	FindByName(name string) (*domain.Medicine, error)
	FindByCode(code string) (*domain.Medicine, error)
}

type MedicineRepositoryPostgresDB interface {
	Query(dataModel *models.MedicinesQueryRequestParams) ([]*domain.Medicine, error)
	Count(dataModel *models.MedicinesQueryRequestParams) (int64, error)
	Create(payload *models.MedicinesCreateRequestBody) (*domain.Medicine, error)
	Delete(payload *models.MedicinesDeleteRequestBody) ([]int64, error)
	Update(payload *models.MedicinesCreateRequestBody, id int64) (*domain.Medicine, error)
	GetMedicinesByIds(ids []int64) ([]*domain.Medicine, error)
}
