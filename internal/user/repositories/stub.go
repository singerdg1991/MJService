package repositories

import (
	"errors"
	"fmt"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"time"

	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/models"
)

type UserRepositoryStub struct {
	Users []*domain.User
}

type userTestCondition struct {
	HasError bool
}

var UserTestCondition *userTestCondition = &userTestCondition{}

func NewUserRepositoryStub() *UserRepositoryStub {
	return &UserRepositoryStub{
		Users: []*domain.User{
			{
				ID:           1,
				FirstName:    "Mohammad Hossein",
				LastName:     "Taher",
				Username:     "saeed",
				Email:        "saeed@gmail.com",
				Phone:        "+989109649362",
				NationalCode: "111111111",
				BirthDate:    time.Now(),
				AvatarUrl:    "https://picsum.photos/100",
				SuspendedAt:  nil,
			},
			{
				ID:           2,
				FirstName:    "Saeed",
				LastName:     "Ghanbari",
				Username:     "sgh370",
				Email:        "sgh370@yahoo.com",
				Phone:        "+989034005707",
				NationalCode: "111111111",
				BirthDate:    time.Now(),
				AvatarUrl:    "https://picsum.photos/100",
				SuspendedAt:  nil,
			},
			{
				ID:           3,
				FirstName:    "Milad",
				LastName:     "Mizani",
				Username:     "milad",
				Email:        "miladmizani@yahoo.com",
				Phone:        "+358442362057",
				NationalCode: "111111111",
				BirthDate:    time.Now(),
				AvatarUrl:    "https://picsum.photos/100",
				SuspendedAt:  nil,
			},
		},
	}
}

func (r *UserRepositoryStub) Query(dataModel *models.UsersQueryRequestParams) ([]*domain.User, error) {
	var users []*domain.User
	for _, v := range r.Users {
		if v.ID == uint(dataModel.ID) ||
			v.FirstName == fmt.Sprintf("%v", dataModel.Filters.FirstName) ||
			v.LastName == fmt.Sprintf("%v", dataModel.Filters.LastName) ||
			v.Username == fmt.Sprintf("%v", dataModel.Filters.Username) ||
			v.Email == fmt.Sprintf("%v", dataModel.Filters.Email) ||
			v.Phone == fmt.Sprintf("%v", dataModel.Filters.Phone) ||
			v.NationalCode == fmt.Sprintf("%v", dataModel.Filters.NationalCode) ||
			v.AvatarUrl == fmt.Sprintf("%v", dataModel.Filters.AvatarUrl) {
			users = append(users, v)
			break
		}
	}
	return users, nil
}

func (r *UserRepositoryStub) Count(dataModel *models.UsersQueryRequestParams) (int64, error) {
	var users []*domain.User
	for _, v := range r.Users {
		if v.ID == uint(dataModel.ID) ||
			v.FirstName == fmt.Sprintf("%v", dataModel.Filters.FirstName) ||
			v.LastName == fmt.Sprintf("%v", dataModel.Filters.LastName) ||
			v.Username == fmt.Sprintf("%v", dataModel.Filters.Username) ||
			v.Email == fmt.Sprintf("%v", dataModel.Filters.Email) ||
			v.Phone == fmt.Sprintf("%v", dataModel.Filters.Phone) ||
			v.NationalCode == fmt.Sprintf("%v", dataModel.Filters.NationalCode) ||
			v.AvatarUrl == fmt.Sprintf("%v", dataModel.Filters.AvatarUrl) {
			users = append(users, v)
			break
		}
	}
	return int64(len(users)), nil
}

func (r *UserRepositoryStub) Create(data *models.UsersCreateRequestBody) (*domain.User, error) {
	user := domain.User{
		ID:           uint(len(r.Users) + 1),
		FirstName:    "Milad",
		LastName:     "Mizani",
		Username:     "milad",
		Email:        "miladmizani@yahoo.com",
		Phone:        "+358442362057",
		NationalCode: "111111111",
		BirthDate:    time.Now(),
		AvatarUrl:    "https://picsum.photos/100",
		SuspendedAt:  nil,
	}
	r.Users = append(r.Users, &user)
	return &user, nil
}

func (r *UserRepositoryStub) Update(payload *models.UsersUpdateRequestBody, id int) (*domain.User, error) {
	var user *domain.User
	for _, v := range r.Users {
		if v.ID == uint(id) {
			user = v
			break
		}
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepositoryStub) Delete(data *models.UsersDeleteRequestBody) ([]uint, error) {
	var deletedIds []uint
	for _, v := range data.IDs {
		for i, user := range r.Users {
			if user.ID == uint(v) {
				deletedIds = append(deletedIds, user.ID)
				r.Users = append(r.Users[:i], r.Users[i+1:]...)
				break
			}
		}
	}
	return deletedIds, nil
}

func (r *UserRepositoryStub) Migrate() {
	// do stuff
}

func (r *UserRepositoryStub) Seed(users []*domain.User) {
	// do stuff
}

func (r *UserRepositoryStub) UpdateAcceptPolicy(acceptPolicy bool, id int) (*domain.User, error) {
	// do stuff
	return nil, nil
}

func (r *UserRepositoryStub) UpdatePassword(newPassword string, id int) (*domain.User, error) {
	// do stuff
	return nil, nil
}

func (r *UserRepositoryStub) CreateUserForCustomer(payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepositoryStub) DeleteUserForCustomer(userID int64) error {
	return nil
}

func (r *UserRepositoryStub) UpdateCustomerIDForUser(userID int64, customerID int64) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepositoryStub) UpdateUserAdditionalInfoForCustomer(payload *sharedmodels.UpdateUserAdditionalInfoForCustomer) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepositoryStub) UpdateUserForCustomer(userID int64, payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error) {
	return nil, nil
}
