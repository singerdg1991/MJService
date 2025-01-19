package service

import (
	"context"
	"errors"
	"github.com/hoitek/Go-Quilder/filters"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Kit/restypes"
	"github.com/hoitek/Maja-Service/internal/_shared/minio"
	"github.com/hoitek/Maja-Service/internal/_shared/sharedmodels"
	"github.com/hoitek/Maja-Service/internal/user/constants"
	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/models"
	"github.com/hoitek/Maja-Service/internal/user/ports"
	"github.com/hoitek/Maja-Service/storage"
	"math"
	"strings"
)

type ContextKey string

type UserService struct {
	PostgresRepository ports.UserRepositoryPostgresDB
	MongoDBRepository  ports.UserRepositoryMongoDB
	MinIOStorage       *storage.MinIO
}

func NewUserService(pDB ports.UserRepositoryPostgresDB, mDB ports.UserRepositoryMongoDB, m *storage.MinIO) ports.UserService {
	go minio.SetupMinIOStorage(constants.USER_BUCKET_NAME, m)
	return &UserService{
		PostgresRepository: pDB,
		MongoDBRepository:  mDB,
		MinIOStorage:       m,
	}
}

func (s *UserService) BindUserToContext(ctx context.Context, user *domain.User) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ContextKey("user"), user)
}

func (s *UserService) BindTokenToContext(ctx context.Context, token string) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ContextKey("token"), token)
}

func (s *UserService) GetUserFromContext(ctx context.Context) *domain.User {
	if ctx == nil {
		return nil
	}
	userEntity, ok := ctx.Value(ContextKey("user")).(*domain.User)
	if !ok {
		return nil
	}
	user, err := s.FindByID(int(userEntity.ID))
	if err != nil || user == nil {
		return nil
	}
	return user
}

func (s *UserService) GetTokenFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	tokenInterface := ctx.Value(ContextKey("token"))
	if tokenInterface == nil {
		return ""
	}
	token, ok := tokenInterface.(string)
	if !ok {
		return ""
	}
	return token
}

func (s *UserService) Query(q *models.UsersQueryRequestParams) (*restypes.QueryResponse, error) {
	q.Page = exp.TerIf(q.Page < 1, 1, q.Page)
	q.Limit = exp.TerIf(q.Limit < 1, 1000, q.Limit)

	users, err := s.PostgresRepository.Query(q)
	if err != nil {
		return nil, err
	}

	count, err := s.PostgresRepository.Count(q)
	if err != nil {
		return nil, err
	}

	page := q.Page
	limit := q.Limit
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 && count > 0 {
		totalPages = page
	}

	return &restypes.QueryResponse{
		Items:      users,
		Limit:      limit,
		Offset:     offset,
		Page:       page,
		TotalRows:  count,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) FindByEmail(email string) (*domain.User, error) {
	queries := &models.UsersQueryRequestParams{
		Filters: models.UserFilterType{
			Email: filters.FilterValue[string]{
				Value: email,
				Op:    operators.EQUALS,
			},
		},
	}
	users, err := s.PostgresRepository.Query(queries)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func (s *UserService) Create(payload *models.UsersCreateRequestBody) (*domain.User, error) {
	return s.PostgresRepository.Create(payload)
}

func (s *UserService) Update(payload *models.UsersUpdateRequestBody, id int) (*domain.User, error) {
	return s.PostgresRepository.Update(payload, id)
}

func (s *UserService) UpdateAcceptPolicy(acceptPolicy bool, id int) (*domain.User, error) {
	return s.PostgresRepository.UpdateAcceptPolicy(acceptPolicy, id)
}

func (s *UserService) Delete(data *models.UsersDeleteRequestBody) (*restypes.DeleteResponse, error) {
	ids, err := s.PostgresRepository.Delete(data)
	if err != nil {
		return nil, err
	}
	return &restypes.DeleteResponse{
		IDs: ids,
	}, nil
}

func (s *UserService) UpdatePassword(newPassword string, id int) (*domain.User, error) {
	return s.PostgresRepository.UpdatePassword(newPassword, id)
}

func (s *UserService) FindByID(id int) (*domain.User, error) {
	queries := &models.UsersQueryRequestParams{
		ID: id,
	}
	users, err := s.PostgresRepository.Query(queries)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func (s *UserService) AssertToUserDomain(payload interface{}) *domain.User {
	result, ok := payload.(*domain.User)
	if ok {
		return result
	}
	return nil
}

func (s *UserService) FindByWorkPhoneNumber(workPhoneNumber string) (*domain.User, error) {
	queries := &models.UsersQueryRequestParams{
		Filters: models.UserFilterType{
			WorkPhoneNumber: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: workPhoneNumber,
			},
		},
	}

	users, err := s.PostgresRepository.Query(queries)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func (s *UserService) FindByUsername(username string) (*domain.User, error) {
	queries := &models.UsersQueryRequestParams{
		Filters: models.UserFilterType{
			Username: filters.FilterValue[string]{
				Op:    operators.EQUALS,
				Value: username,
			},
		},
	}

	users, err := s.PostgresRepository.Query(queries)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func (s *UserService) CreateUserForCustomer(payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error) {
	return s.PostgresRepository.CreateUserForCustomer(payload)
}

func (s *UserService) DeleteUserForCustomer(userID int64) error {
	return s.PostgresRepository.DeleteUserForCustomer(userID)
}

func (s *UserService) UpdateCustomerIDForUser(userID int64, customerID int64) (*domain.User, error) {
	return s.PostgresRepository.UpdateCustomerIDForUser(userID, customerID)
}

func (s *UserService) UpdateUserAdditionalInfoForCustomer(payload *sharedmodels.UpdateUserAdditionalInfoForCustomer) (*domain.User, error) {
	return s.PostgresRepository.UpdateUserAdditionalInfoForCustomer(payload)
}

func (s *UserService) UpdateUserForCustomer(userID int64, payload *sharedmodels.CustomersCreatePersonalInfo) (*domain.User, error) {
	return s.PostgresRepository.UpdateUserForCustomer(userID, payload)
}

func (s *UserService) IsDispatcher(userID int64) bool {
	users, err := s.PostgresRepository.Query(&models.UsersQueryRequestParams{
		ID: int(userID),
	})
	if err != nil {
		return false
	}
	if len(users) < 1 {
		return false
	}
	user := users[0]
	for _, role := range user.Roles {
		if strings.ToLower(role.Name) == "dispatcher" {
			return true
		}
	}
	return false
}
