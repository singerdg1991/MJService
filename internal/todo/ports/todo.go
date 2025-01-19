package ports

import (
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/todo/domain"
	"github.com/hoitek/Maja-Service/internal/todo/models"
)

type TodoService interface {
	Query(dataModel *models.TodosQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.TodosCreateRequestBody) (*domain.Todo, error)
	Delete(payload *models.TodosDeleteRequestBody) (*restypes.DeleteResponse, error)
	Update(payload *models.TodosCreateRequestBody, id int64) (*domain.Todo, error)
	GetTodosByIds(ids []int64) ([]*domain.Todo, error)
	FindByID(id int64) (*domain.Todo, error)
	FindByTitle(title string) (*domain.Todo, error)
	UpdateStatus(payload *models.TodosUpdateStatusRequestBody, id int64) (*domain.Todo, error)
}

type TodoRepositoryPostgresDB interface {
	Query(dataModel *models.TodosQueryRequestParams) ([]*domain.Todo, error)
	Count(dataModel *models.TodosQueryRequestParams) (int64, error)
	Create(payload *models.TodosCreateRequestBody) (*domain.Todo, error)
	Delete(payload *models.TodosDeleteRequestBody) ([]int64, error)
	Update(payload *models.TodosCreateRequestBody, id int64) (*domain.Todo, error)
	GetTodosByIds(ids []int64) ([]*domain.Todo, error)
	UpdateStatus(payload *models.TodosUpdateStatusRequestBody, id int64) (*domain.Todo, error)
}
