package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/todo/domain"
	"github.com/hoitek/Maja-Service/internal/todo/models"
)

type TodoRepositoryStub struct {
	Todos []*domain.Todo
}

type todoTestCondition struct {
	HasError bool
}

var UserTestCondition *todoTestCondition = &todoTestCondition{}

func NewTodoRepositoryStub() *TodoRepositoryStub {
	return &TodoRepositoryStub{
		Todos: []*domain.Todo{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *TodoRepositoryStub) Query(dataModel *models.TodosQueryRequestParams) ([]*domain.Todo, error) {
	var todos []*domain.Todo
	for _, v := range r.Todos {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			todos = append(todos, v)
			break
		}
	}
	return todos, nil
}

func (r *TodoRepositoryStub) Count(dataModel *models.TodosQueryRequestParams) (int64, error) {
	var todos []*domain.Todo
	for _, v := range r.Todos {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			todos = append(todos, v)
			break
		}
	}
	return int64(len(todos)), nil
}

func (r *TodoRepositoryStub) Migrate() {
	// do stuff
}

func (r *TodoRepositoryStub) Seed() {
	// do stuff
}

func (r *TodoRepositoryStub) Create(payload *models.TodosCreateRequestBody) (*domain.Todo, error) {
	panic("implement me")
}

func (r *TodoRepositoryStub) Delete(payload *models.TodosDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *TodoRepositoryStub) Update(payload *models.TodosCreateRequestBody, id int64) (*domain.Todo, error) {
	panic("implement me")
}

func (r *TodoRepositoryStub) GetTodosByIds(ids []int64) ([]*domain.Todo, error) {
	panic("implement me")
}

func (r *TodoRepositoryStub) UpdateStatus(payload *models.TodosUpdateStatusRequestBody, id int64) (*domain.Todo, error) {
	panic("implement me")
}
