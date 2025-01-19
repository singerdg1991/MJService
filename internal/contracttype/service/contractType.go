package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/contracttype/constants"
	"github.com/hoitek/Maja-Service/internal/contracttype/domain"
	"github.com/hoitek/Maja-Service/internal/contracttype/models"
	"github.com/hoitek/Maja-Service/internal/contracttype/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type ContractTypeService struct {
	PostgresRepository ports.ContractTypeRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewContractTypeService(pDB ports.ContractTypeRepositoryPostgresDB, m *storage.MinIO) ContractTypeService {
	go minio.SetupMinIOStorage(constants.CONTRACT_TYPE_BUCKET_NAME, m)
	return ContractTypeService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *ContractTypeService) Query(q *models.ContractTypesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying ContractTypes", q)
	contractTypes, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.ContractTypesQueryRequestParams{
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
	var items []*domain.ContractType
	for _, item := range contractTypes {
		items = append(items, &domain.ContractType{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
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

func (s *ContractTypeService) Create(payload *models.ContractTypesCreateRequestBody) (*domain.ContractType, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *ContractTypeService) Delete(payload *models.ContractTypesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *ContractTypeService) Update(payload *models.ContractTypesCreateRequestBody, id int64) (*domain.ContractType, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *ContractTypeService) GetContractTypesByIds(ids []int64) ([]*domain.ContractType, error) {
	return s.PostgresRepository.GetContractTypesByIds(ids)
}

func (s *ContractTypeService) FindByID(id int64) (*domain.ContractType, error) {
	r, err := s.Query(&models.ContractTypesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	contractTypes := r.Items.([]*domain.ContractType)
	return contractTypes[0], nil
}

func (s *ContractTypeService) FindByName(name string) (*domain.ContractType, error) {
	r, err := s.Query(&models.ContractTypesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ContractTypeFilterType{
			Name: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("language skill not found")
	}
	contractTypes := r.Items.([]*domain.ContractType)
	return contractTypes[0], nil
}
