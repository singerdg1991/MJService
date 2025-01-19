package ports

import (
	"context"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"

	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/models"
)

type UserService interface {
	Query(q *models.UsersQueryRequestParams) (*restypes.QueryResponse, error)
	Create(payload *models.UsersCreateRequestBody) (*domain.User, error)
	Update(payload *models.UsersUpdateRequestBody, id int) (*domain.User, error)
	UpdateAcceptPolicy(acceptPolicy bool, id int) (*domain.User, error)
	UpdatePassword(newPassword string, id int) (*domain.User, error)
	Delete(payload *models.UsersDeleteRequestBody) (*restypes.DeleteResponse, error)
	BindUserToContext(ctx context.Context, user *domain.User) context.Context
	BindTokenToContext(ctx context.Context, token string) context.Context
	GetUserFromContext(ctx context.Context) *domain.User
	GetTokenFromContext(ctx context.Context) string
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
	AssertToUserDomain(payload interface{}) *domain.User
	FindByWorkPhoneNumber(workPhoneNumber string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	CreateUserForCustomer(payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error)
	UpdateUserForCustomer(userID int64, payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error)
	DeleteUserForCustomer(userID int64) error
	UpdateCustomerIDForUser(userID int64, customerID int64) (*domain.User, error)
	UpdateUserAdditionalInfoForCustomer(payload *sharedmodels.UpdateUserAdditionalInfoForCustomer) (*domain.User, error)
	IsDispatcher(userID int64) bool
}

type UserRepositoryPostgresDB interface {
	Query(dataModel *models.UsersQueryRequestParams) ([]*domain.User, error)
	Count(dataModel *models.UsersQueryRequestParams) (int64, error)
	Create(payload *models.UsersCreateRequestBody) (*domain.User, error)
	Update(payload *models.UsersUpdateRequestBody, id int) (*domain.User, error)
	UpdateAcceptPolicy(acceptPolicy bool, id int) (*domain.User, error)
	UpdatePassword(newPassword string, id int) (*domain.User, error)
	Delete(payload *models.UsersDeleteRequestBody) ([]uint, error)
	CreateUserForCustomer(payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error)
	UpdateUserForCustomer(userID int64, payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error)
	DeleteUserForCustomer(userID int64) error
	UpdateCustomerIDForUser(userID int64, customerID int64) (*domain.User, error)
	UpdateUserAdditionalInfoForCustomer(payload *sharedmodels.UpdateUserAdditionalInfoForCustomer) (*domain.User, error)
}

type UserRepositoryMongoDB interface {
	Query(queries *models.UsersQueryRequestParams) ([]*domain.User, error)
	Count(queries *models.UsersQueryRequestParams) (int64, error)
	Create(payload map[string]interface{}) (interface{}, error)
	Update(payload *domain.User, id int) error
	Delete(ids []uint) error
	UpdateByPostgresID(postgresID int, payload interface{}) (interface{}, error)
}
