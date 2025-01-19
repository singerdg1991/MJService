package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/prescription/domain"
	"github.com/hoitek/Maja-Service/internal/prescription/models"
)

type PrescriptionService interface {
	Query(dataModel *models.PrescriptionsQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.PrescriptionsCreateRequestBody) (*domain.Prescription, error)
	Delete(payload *models.PrescriptionsDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.PrescriptionsUpdateRequestBody, id int) (*domain.Prescription, error)
	GetPrescriptionByID(id int64) (*domain.Prescription, error)
	UpdatePrescriptionAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Prescription, error)
}

type PrescriptionRepositoryPostgresDB interface {
	Query(dataModel *models.PrescriptionsQueryRequestParams) ([]*domain.Prescription, error)
	Count(dataModel *models.PrescriptionsQueryRequestParams) (int64, error)
	Create(payload *models.PrescriptionsCreateRequestBody) (*domain.Prescription, error)
	Delete(payload *models.PrescriptionsDeleteRequestBody) ([]int64, error)
	Update(payload *models.PrescriptionsUpdateRequestBody, id int) (*domain.Prescription, error)
	UpdatePrescriptionAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Prescription, error)
}
