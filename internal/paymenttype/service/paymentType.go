package service

import (
	"errors"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/paymenttype/constants"
	"github.com/hoitek/Maja-Service/internal/paymenttype/domain"
	"github.com/hoitek/Maja-Service/internal/paymenttype/models"
	"github.com/hoitek/Maja-Service/internal/paymenttype/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type PaymentTypeService struct {
	PostgresRepository ports.PaymentTypeRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewPaymentTypeService(pDB ports.PaymentTypeRepositoryPostgresDB, m *storage.MinIO) PaymentTypeService {
	go minio.SetupMinIOStorage(constants.PAYMENT_TYPE_BUCKET_NAME, m)
	return PaymentTypeService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *PaymentTypeService) Query(q *models.PaymentTypesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying paymentTypes", q)
	paymentTypes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.PaymentTypesQueryRequestParams{
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

	// Transform the response to the format that the frontend expects
	var items []*domain.PaymentType
	for _, item := range paymentTypes {
		items = append(items, &domain.PaymentType{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return &restypes.QueryResponse{
		Items:      items,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *PaymentTypeService) GetPaymentTypeByID(id int) (*domain.PaymentType, error) {
	r, err := s.Query(&models.PaymentTypesQueryRequestParams{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows <= 0 {
		return nil, err
	}
	paymentTypes, ok := r.Items.([]*domain.PaymentType)
	if !ok {
		return nil, errors.New("paymentType not found")
	}
	return paymentTypes[0], nil
}
