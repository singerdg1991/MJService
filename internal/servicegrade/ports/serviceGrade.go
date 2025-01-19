package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/servicegrade/domain"
	"github.com/hoitek/Maja-Service/internal/servicegrade/models"
)

type ServiceGradeService interface {
	Query(dataModel *models.ServiceGradesQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.ServiceGradesCreateRequestBody) (*domain.ServiceGrade, error)
	Delete(payload *models.ServiceGradesDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.ServiceGradesCreateRequestBody, id int64) (*domain.ServiceGrade, error)
	GetServiceGradesByIds(ids []int64) ([]*domain.ServiceGrade, error)
	FindByID(id int64) (*domain.ServiceGrade, error)
	FindByName(name string) (*domain.ServiceGrade, error)
}

type ServiceGradeRepositoryPostgresDB interface {
	Query(dataModel *models.ServiceGradesQueryRequestParams) ([]*domain.ServiceGrade, error)
	Count(dataModel *models.ServiceGradesQueryRequestParams) (int64, error)
	Create(payload *models.ServiceGradesCreateRequestBody) (*domain.ServiceGrade, error)
	Delete(payload *models.ServiceGradesDeleteRequestBody) ([]int64, error)
	Update(payload *models.ServiceGradesCreateRequestBody, id int64) (*domain.ServiceGrade, error)
	GetServiceGradesByIds(ids []int64) ([]*domain.ServiceGrade, error)
}
