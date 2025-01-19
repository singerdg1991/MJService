package service

import (
	"errors"
	"log"
	"math"

	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/todo/constants"
	"github.com/hoitek/Maja-Service/internal/todo/domain"
	"github.com/hoitek/Maja-Service/internal/todo/models"
	"github.com/hoitek/Maja-Service/internal/todo/ports"
	"github.com/hoitek/Maja-Service/storage"

	"github.com/hoitek/Kit/restypes"

	"github.com/hoitek/Kit/exp"
)

type TodoService struct {
	PostgresRepository ports.TodoRepositoryPostgresDB
	MinIOStorage       *storage.MinIO
}

func NewTodoService(pDB ports.TodoRepositoryPostgresDB, m *storage.MinIO) TodoService {
	go minio.SetupMinIOStorage(constants.TODO_BUCKET_NAME, m)
	return TodoService{
		PostgresRepository: pDB,
		MinIOStorage:       m,
	}
}

func (s *TodoService) Query(q *models.TodosQueryRequestParams) (*restypes.QueryResponse, error) {
	log.Println("Querying todos", q)
	todos, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(&models.TodosQueryRequestParams{
		ID:      q.ID,
		UserID:  q.UserID,
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
	var items []*domain.Todo
	for _, item := range todos {
		items = append(items, &domain.Todo{
			ID:            item.ID,
			UserID:        item.UserID,
			Title:         item.Title,
			Date:          item.Date,
			Time:          item.Time,
			DateStr:       item.DateStr,
			TimeStr:       item.TimeStr,
			User:          item.User,
			Description:   item.Description,
			Status:        item.Status,
			DoneAt:        item.DoneAt,
			CreatedByUser: item.CreatedByUser,
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

func (s *TodoService) Create(payload *models.TodosCreateRequestBody) (*domain.Todo, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *TodoService) Delete(payload *models.TodosDeleteRequestBody) (*restypes.DeleteResponse, error) {
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

func (s *TodoService) Update(payload *models.TodosCreateRequestBody, id int64) (*domain.Todo, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *TodoService) GetTodosByIds(ids []int64) ([]*domain.Todo, error) {
	return s.PostgresRepository.GetTodosByIds(ids)
}

func (s *TodoService) FindByID(id int64) (*domain.Todo, error) {
	r, err := s.Query(&models.TodosQueryRequestParams{
		ID: int(id),
	})
	if err != nil {
		return nil, err
	}
	if r.TotalRows == 0 {
		return nil, errors.New("todo not found")
	}
	todos := r.Items.([]*domain.Todo)
	return todos[0], nil
}

func (s *TodoService) FindByTitle(title string) (*domain.Todo, error) {
	r, err := s.Query(&models.TodosQueryRequestParams{
		Page:  1,
		Limit: 1,
		Filters: models.TodoFilterType{
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
		return nil, errors.New("todo not found")
	}
	todos, ok := r.Items.([]*domain.Todo)
	if !ok {
		return nil, errors.New("todo not found")
	}
	if len(todos) < 1 {
		return nil, errors.New("todo not found")
	}
	return todos[0], nil
}

func (s *TodoService) UpdateStatus(payload *models.TodosUpdateStatusRequestBody, id int64) (*domain.Todo, error) {
	return s.PostgresRepository.UpdateStatus(payload, id)
}
