package repositories

import (
	"fmt"

	"github.com/hoitek/Maja-Service/internal/email/domain"
	"github.com/hoitek/Maja-Service/internal/email/models"
)

type EmailRepositoryStub struct {
	Emails []*domain.Email
}

type emailTestCondition struct {
	HasError bool
}

var UserTestCondition *emailTestCondition = &emailTestCondition{}

func NewEmailRepositoryStub() *EmailRepositoryStub {
	return &EmailRepositoryStub{
		Emails: []*domain.Email{
			{
				ID:    1,
				Title: "test",
			},
		},
	}
}

func (r *EmailRepositoryStub) Query(dataModel *models.EmailsQueryRequestParams) ([]*domain.Email, error) {
	var emails []*domain.Email
	for _, v := range r.Emails {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			emails = append(emails, v)
			break
		}
	}
	return emails, nil
}

func (r *EmailRepositoryStub) Count(dataModel *models.EmailsQueryRequestParams) (int64, error) {
	var emails []*domain.Email
	for _, v := range r.Emails {
		if v.ID == uint(dataModel.ID) ||
			v.Title == fmt.Sprintf("%v", dataModel.Filters.Title) {
			emails = append(emails, v)
			break
		}
	}
	return int64(len(emails)), nil
}

func (r *EmailRepositoryStub) Migrate() {
	// do stuff
}

func (r *EmailRepositoryStub) Seed() {
	// do stuff
}

func (r *EmailRepositoryStub) Create(payload *models.EmailsCreateRequestBody) (*domain.Email, error) {
	panic("implement me")
}

func (r *EmailRepositoryStub) Delete(payload *models.EmailsDeleteRequestBody) ([]int64, error) {
	panic("implement me")
}

func (r *EmailRepositoryStub) UpdateCategory(payload *models.EmailsUpdateCategoryRequestBody, id int64) (*domain.Email, error) {
	panic("implement me")
}

func (r *EmailRepositoryStub) UpdateStar(payload *models.EmailsUpdateStarRequestBody, id int64) (*domain.Email, error) {
	panic("implement me")
}

func (r *EmailRepositoryStub) GetEmailsByIds(ids []int64) ([]*domain.Email, error) {
	panic("implement me")
}
