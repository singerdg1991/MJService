package service

import (
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/_shared/types"
	"github.com/hoitek/Maja-Service/internal/archive/constants"
	"github.com/hoitek/Maja-Service/internal/archive/domain"
	"github.com/hoitek/Maja-Service/internal/archive/models"
	"github.com/hoitek/Maja-Service/internal/archive/ports"
	"github.com/hoitek/Maja-Service/storage"
	"log"
	"math"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type ArchiveService struct {
	PostgresRepository ports.ArchiveRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewArchiveService(pDB ports.ArchiveRepositoryPostgresDB, m *storage.MinIO) ArchiveService {
	go minio.SetupMinIOStorage(constants.ARCHIVE_BUCKET_NAME, m)
	return ArchiveService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *ArchiveService) Query(q *models.ArchivesQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying archives", q)
	archives, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.ArchivesQueryRequestParams{
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
	var items []*domain.Archive
	for _, item := range archives {
		items = append(items, &domain.Archive{
			ID:          item.ID,
			UserID:      item.UserID,
			User:        item.User,
			Title:       item.Title,
			Subject:     item.Subject,
			Description: item.Description,
			Attachments: item.Attachments,
			Date:        item.Date,
			Time:        item.Time,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			DeletedAt:   item.DeletedAt,
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

func (s *ArchiveService) Create(payload *models.ArchivesCreateRequestBody) (*domain.Archive, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *ArchiveService) Delete(payload *models.ArchivesDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *ArchiveService) Update(payload *models.ArchivesCreateRequestBody, id int64) (*domain.Archive, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *ArchiveService) GetArchivesByIds(ids []int64) ([]*domain.Archive, error) {
	return s.PostgresRepository.GetArchivesByIds(ids)
}

func (s *ArchiveService) FindByID(id int64) (*domain.Archive, error) {
	r, err := s.Query(&models.ArchivesQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("archive not found")
	}
	archives := r.Items.([]*domain.Archive)
	return archives[0], nil
}

func (s *ArchiveService) FindByTitle(title string) (*domain.Archive, error) {
	r, err := s.Query(&models.ArchivesQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.ArchiveFilterType{
			Title: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: title,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("archive not found")
	}
	archives := r.Items.([]*domain.Archive)
	return archives[0], nil
}

func (s *ArchiveService) UpdateAttachments(attachments []*types.UploadMetadata, id int64) (*domain.Archive, error) {
	archive, err := s.PostgresRepository.UpdateAttachments(attachments, id)
	if err != nil {
		return nil, err
	}
	return archive, nil
}
