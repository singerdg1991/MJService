package repositories

import (
	"github.com/hoitek/Maja-Service/internal/staff/domain"
	"github.com/hoitek/Maja-Service/internal/staff/models"
)

type StaffRepositoryMongoDBStub struct {
}

func NewStaffRepositoryMongoDBStub() *StaffRepositoryMongoDBStub {
	return &StaffRepositoryMongoDBStub{}
}

func (r *StaffRepositoryMongoDBStub) Query(queries *models.StaffsQueryRequestParams) ([]*domain.Staff, error) {
	return []*domain.Staff{}, nil
}

func (r *StaffRepositoryMongoDBStub) Count(queries *models.StaffsQueryRequestParams) (int64, error) {
	return 0, nil
}

func (r *StaffRepositoryMongoDBStub) Create(postgresID int, payload interface{}) (interface{}, error) {
	return nil, nil
}

func (r *StaffRepositoryMongoDBStub) Update(payload interface{}, id int) error {
	return nil
}

func (r *StaffRepositoryMongoDBStub) CreateOrUpdate(postgresID int, payload interface{}) (interface{}, error) {
	return nil, nil
}

func (r *StaffRepositoryMongoDBStub) Delete(ids []uint) error {
	return nil
}

// UpdateByPostgresID updates doc by postgres id
func (r *StaffRepositoryMongoDBStub) UpdateByPostgresID(postgresID int, payload interface{}) (interface{}, error) {
	return nil, nil
}

// UpdateUserInfo updates user info
func (r *StaffRepositoryMongoDBStub) UpdateUserInfo(userID int, payload interface{}) (interface{}, error) {
	return nil, nil
}

// CreateEmptyForUserID creates staff
func (r *StaffRepositoryMongoDBStub) CreateEmptyForUserID(payload interface{}) (interface{}, error) {
	return nil, nil
}
