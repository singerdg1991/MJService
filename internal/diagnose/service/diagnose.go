package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/diagnose/constants"
	"github.com/hoitek/Maja-Service/internal/diagnose/domain"
	"github.com/hoitek/Maja-Service/internal/diagnose/models"
	"github.com/hoitek/Maja-Service/internal/diagnose/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type DiagnoseService struct {
	PostgresRepository ports.DiagnoseRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewDiagnoseService(pDB ports.DiagnoseRepositoryPostgresDB, m *storage.MinIO) DiagnoseService {
	go minio.SetupMinIOStorage(constants.DIAGNOSE_BUCKET_NAME, m)
	return DiagnoseService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *DiagnoseService) Query(q *models.DiagnosesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying diagnoses", q)
	diagnoses, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.DiagnosesQueryRequestParams{
		ID:      q.ID,
		Page:    q.Page,
		Limit:   0,
		Filters: q.Filters,
	})
	if err != nil {
		return nil, err
	}

	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 10, 1, q.Limit)

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      diagnoses,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *DiagnoseService) Create(payload *models.DiagnosesCreateRequestBody) (*domain.Diagnose, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *DiagnoseService) Delete(payload *models.DiagnosesDeleteRequestBody) (*restypes.DeleteResponse, error) {
	deletedIds, err := s.PostgresRepository.Delete(payload)
	if err != nil {
		return nil, err
	}

	// TODO this is a temporary solution, we should return the deleted ids as int64 we show change restypes.DeleteResponse.IDs to []int64
	var ids []uint
	for _, id := range deletedIds {
		ids = append(ids, uint(id))
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *DiagnoseService) Update(payload *models.DiagnosesCreateRequestBody, id int) (*domain.Diagnose, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *DiagnoseService) GetDiagnoseByID(id int64) (*domain.Diagnose, error) {
	diagnoses, err := s.PostgresRepository.Query(&models.DiagnosesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(diagnoses) == 0 {
		return nil, nil
	}
	return diagnoses[0], nil
}
