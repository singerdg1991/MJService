package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/keikkala/constants"
	"github.com/hoitek/Maja-Service/internal/keikkala/domain"
	"github.com/hoitek/Maja-Service/internal/keikkala/models"
	"github.com/hoitek/Maja-Service/internal/keikkala/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type KeikkalaService struct {
	PostgresRepository ports.KeikkalaRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewKeikkalaService(pDB ports.KeikkalaRepositoryPostgresDB, m *storage.MinIO) KeikkalaService {
	go minio.SetupMinIOStorage(constants.KEIKKALA_BUCKET_NAME, m)
	return KeikkalaService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *KeikkalaService) Query(q *models.KeikkalasQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying keikkalas", q)
	keikkalas, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.KeikkalasQueryRequestParams{
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
		Items:      keikkalas,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *KeikkalaService) Create(payload *models.KeikkalasCreateRequestBody) (*domain.Keikkala, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *KeikkalaService) Delete(payload *models.KeikkalasDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *KeikkalaService) GetKeikkalaShiftsByIds(ids []int64) ([]*domain.Keikkala, error) {
	return s.PostgresRepository.GetKeikkalaShiftsByIds(ids)
}

func (s *KeikkalaService) FindByID(id int64) (*domain.Keikkala, error) {
	r, err := s.Query(&models.KeikkalasQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("keikkala not found")
	}
	keikkalas := r.Items.([]*domain.Keikkala)
	return keikkalas[0], nil
}

func (s *KeikkalaService) QueryShiftStatistics(queries *models.KeikkalasQueryShiftStatisticsRequestParams) (*models.KeikkalasQueryShiftStatisticsResponseData, error) {
	morningCount, eveningCount, nightCount, err := s.PostgresRepository.QueryShiftStatistics(queries)
	if err != nil {
		return nil, err
	}

	return &models.KeikkalasQueryShiftStatisticsResponseData{
		MorningCount: morningCount,
		EveningCount: eveningCount,
		NightCount:   nightCount,
	}, nil
}
