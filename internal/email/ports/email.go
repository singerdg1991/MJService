package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/email/domain"
	"github.com/hoitek/Maja-Service/internal/email/models"
)

type EmailService interface {
	Query(dataModel *models.EmailsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.EmailsCreateRequestBody) (*domain.Email, error)
	Delete(payload *models.EmailsDeleteRequestBody) (*restypes.DeleteResponse, error)
	UpdateCategory(payload *models.EmailsUpdateCategoryRequestBody, id int64) (*domain.Email, error)
	UpdateStar(payload *models.EmailsUpdateStarRequestBody, id int64) (*domain.Email, error)
	GetEmailsByIds(ids []int64) ([]*domain.Email, error)
	FindByID(id int64) (*domain.Email, error)
}

type EmailRepositoryPostgresDB interface {
	Query(dataModel *models.EmailsQueryRequestParams) ([]*domain.Email, error)
	Count(dataModel *models.EmailsQueryRequestParams) (int64, error)
	Create(payload *models.EmailsCreateRequestBody) (*domain.Email, error)
	Delete(payload *models.EmailsDeleteRequestBody) ([]int64, error)
	UpdateCategory(payload *models.EmailsUpdateCategoryRequestBody, id int64) (*domain.Email, error)
	UpdateStar(payload *models.EmailsUpdateStarRequestBody, id int64) (*domain.Email, error)
	GetEmailsByIds(ids []int64) ([]*domain.Email, error)
}
