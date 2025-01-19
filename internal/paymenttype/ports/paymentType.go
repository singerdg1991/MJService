package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/paymenttype/domain"
	"github.com/hoitek/Maja-Service/internal/paymenttype/models"
)

type PaymentTypeService interface {
	Query(dataModel *models.PaymentTypesQueryRequestParams) (*restypes.QueryResponse, error)
	GetPaymentTypeByID(id int) (*domain.PaymentType, error)
}

type PaymentTypeRepositoryPostgresDB interface {
	Query(dataModel *models.PaymentTypesQueryRequestParams) ([]*domain.PaymentType, error)
	Count(dataModel *models.PaymentTypesQueryRequestParams) (int64, error)
}
