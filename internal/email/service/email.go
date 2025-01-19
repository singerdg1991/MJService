package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/email/constants"
	"github.com/hoitek/Maja-Service/internal/email/domain"
	"github.com/hoitek/Maja-Service/internal/email/models"
	"github.com/hoitek/Maja-Service/internal/email/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type EmailService struct {
	PostgresRepository ports.EmailRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewEmailService(pDB ports.EmailRepositoryPostgresDB, m *storage.MinIO) EmailService {
	go minio.SetupMinIOStorage(constants.EMAIL_BUCKET_NAME, m)
	return EmailService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *EmailService) Query(q *models.EmailsQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying emails", q)
	emails, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.EmailsQueryRequestParams{
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
		Items:      emails,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *EmailService) Create(payload *models.EmailsCreateRequestBody) (*domain.Email, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *EmailService) Delete(payload *models.EmailsDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *EmailService) UpdateCategory(payload *models.EmailsUpdateCategoryRequestBody, id int64) (*domain.Email, error) {
	return s.PostgresRepository.UpdateCategory(payload, id)
}

func (s *EmailService) UpdateStar(payload *models.EmailsUpdateStarRequestBody, id int64) (*domain.Email, error) {
	return s.PostgresRepository.UpdateStar(payload, id)
}

func (s *EmailService) GetEmailsByIds(ids []int64) ([]*domain.Email, error) {
	return s.PostgresRepository.GetEmailsByIds(ids)
}

func (s *EmailService) FindByID(id int64) (*domain.Email, error) {
	r, err := s.Query(&models.EmailsQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("email not found")
	}
	emails := r.Items.([]*domain.Email)
	return emails[0], nil
}
