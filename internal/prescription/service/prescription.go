package service

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/prescription/constants"
	"github.com/hoitek/Maja-Service/internal/prescription/domain"
	"github.com/hoitek/Maja-Service/internal/prescription/models"
	"github.com/hoitek/Maja-Service/internal/prescription/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type PrescriptionService struct {
	PostgresRepository ports.PrescriptionRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewPrescriptionService(pDB ports.PrescriptionRepositoryPostgresDB, m *storage.MinIO) PrescriptionService {
	go minio.SetupMinIOStorage(constants.PRESCRIPTION_BUCKET_NAME, m)
	return PrescriptionService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *PrescriptionService) Query(q *models.PrescriptionsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying prescriptions", q)
	prescriptions, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.PrescriptionsQueryRequestParams{
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
		Items:      prescriptions,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *PrescriptionService) Create(payload *models.PrescriptionsCreateRequestBody) (*domain.Prescription, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *PrescriptionService) Delete(payload *models.PrescriptionsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *PrescriptionService) Update(payload *models.PrescriptionsUpdateRequestBody, id int) (*domain.Prescription, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *PrescriptionService) GetPrescriptionByID(id int64) (*domain.Prescription, error) {
	prescriptions, err := s.PostgresRepository.Query(&models.PrescriptionsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if len(prescriptions) == 0 {
		return nil, nil
	}
	return prescriptions[0], nil
}

func (s *PrescriptionService) UpdatePrescriptionAttachments(previousAttachments []types.UploadMetadata, attachments []*types.UploadMetadata, id int64) (*domain.Prescription, error) {
	return s.PostgresRepository.UpdatePrescriptionAttachments(previousAttachments, attachments, id)
}
