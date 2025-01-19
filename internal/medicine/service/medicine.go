package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/medicine/constants"
	"github.com/hoitek/Maja-Service/internal/medicine/domain"
	"github.com/hoitek/Maja-Service/internal/medicine/models"
	"github.com/hoitek/Maja-Service/internal/medicine/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type MedicineService struct {
	PostgresRepository ports.MedicineRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewMedicineService(pDB ports.MedicineRepositoryPostgresDB, m *storage.MinIO) MedicineService {
	go minio.SetupMinIOStorage(constants.MEDICINE_BUCKET_NAME, m)
	return MedicineService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *MedicineService) Query(q *models.MedicinesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying medicines", q)
	medicines, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.MedicinesQueryRequestParams{
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
	var items []*domain.Medicine
	for _, item := range medicines {
		items = append(items, &domain.Medicine{
			ID:           item.ID,
			Name:         item.Name,
			Code:         item.Code,
			Availability: item.Availability,
			Manufacturer: item.Manufacturer,
			PurposeOfUse: item.PurposeOfUse,
			Instruction:  item.Instruction,
			SideEffects:  item.SideEffects,
			Conditions:   item.Conditions,
			Description:  item.Description,
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

func (s *MedicineService) Create(payload *models.MedicinesCreateRequestBody) (*domain.Medicine, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *MedicineService) Delete(payload *models.MedicinesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *MedicineService) Update(payload *models.MedicinesCreateRequestBody, id int64) (*domain.Medicine, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *MedicineService) GetMedicinesByIds(ids []int64) ([]*domain.Medicine, error) {
	return s.PostgresRepository.GetMedicinesByIds(ids)
}

func (s *MedicineService) FindByID(id int64) (*domain.Medicine, error) {
	r, err := s.Query(&models.MedicinesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("medicine not found")
	}
	medicines := r.Items.([]*domain.Medicine)
	return medicines[0], nil
}

func (s *MedicineService) FindByName(name string) (*domain.Medicine, error) {
	r, err := s.Query(&models.MedicinesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.MedicineFilterType{
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
		return nil, errors.New("medicine not found")
	}
	if len(r.Items.([]*domain.Medicine)) == 0 {
		return nil, errors.New("medicine not found")
	}
	medicines := r.Items.([]*domain.Medicine)
	return medicines[0], nil
}

func (s *MedicineService) FindByCode(code string) (*domain.Medicine, error) {
	r, err := s.Query(&models.MedicinesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.MedicineFilterType{
			Code: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: code,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("medicine not found")
	}
	if len(r.Items.([]*domain.Medicine)) == 0 {
		return nil, errors.New("medicine not found")
	}
	medicines := r.Items.([]*domain.Medicine)
	return medicines[0], nil
}
