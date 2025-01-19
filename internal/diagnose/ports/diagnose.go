package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/diagnose/domain"
	"github.com/hoitek/Maja-Service/internal/diagnose/models"
)

type DiagnoseService interface {
	Query(dataModel *models.DiagnosesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.DiagnosesCreateRequestBody) (*domain.Diagnose, error)
	Delete(payload *models.DiagnosesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.DiagnosesCreateRequestBody, id int) (*domain.Diagnose, error)
	GetDiagnoseByID(id int64) (*domain.Diagnose, error)
}

type DiagnoseRepositoryPostgresDB interface {
	Query(dataModel *models.DiagnosesQueryRequestParams) ([]*domain.Diagnose, error)
	Count(dataModel *models.DiagnosesQueryRequestParams) (int64, error)
	Create(payload *models.DiagnosesCreateRequestBody) (*domain.Diagnose, error)
	Delete(payload *models.DiagnosesDeleteRequestBody) ([]int64, error)
	Update(payload *models.DiagnosesCreateRequestBody, id int) (*domain.Diagnose, error)
}
