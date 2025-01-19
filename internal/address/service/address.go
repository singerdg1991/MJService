package service

import (
	"errors"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/address/constants"
	"github.com/hoitek/Maja-Service/internal/address/domain"
	"github.com/hoitek/Maja-Service/internal/address/models"
	"github.com/hoitek/Maja-Service/internal/address/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/exp"
)

type AddressService struct {
	PostgresRepository ports.AddressRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewAddressService(pDB ports.AddressRepositoryPostgresDB, m *storage.MinIO) AddressService {
	go minio.SetupMinIOStorage(constants.ADDRESS_BUCKET_NAME, m)
	return AddressService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *AddressService) Query(q *models.AddressesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying addresses", q)
	addresses, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.AddressesQueryRequestParams{
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
		Items:      addresses,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *AddressService) Create(payload *models.AddressesCreateRequestBody) (*domain.Address, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *AddressService) Delete(payload *models.AddressesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *AddressService) Update(payload *models.AddressesCreateRequestBody, id int) (*domain.Address, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *AddressService) FindByID(id int64) (*domain.Address, error) {
	r, err := s.Query(&models.AddressesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("address not found")
	}
	addresses := r.Items.([]*domain.Address)
	if len(addresses) == 0 {
		return nil, errors.New("address not found")
	}
	return addresses[0], nil
}
