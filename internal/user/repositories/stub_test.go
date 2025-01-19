package repositories

import (
	"testing"

	"github.com/hoitek/Maja-Service/internal/user/models"
)

func TestNewUserRepositoryStub(t *testing.T) {
	repo := NewUserRepositoryStub()

	if len(repo.Users) == 0 {
		t.Error("Users field can not be empty")
	}
}

func TestUserRepositoryStub_Create(t *testing.T) {
	repo := NewUserRepositoryStub()
	_, err := repo.Create(&models.UsersCreateRequestBody{
		FirstName:    "TestName",
		LastName:     "TestLastName",
		Username:     "test",
		Email:        "test@yahoo.com",
		Phone:        "+358442362057",
		NationalCode: "111111111",
		BirthDate:    "1111111",
		AvatarUrl:    "https://picsum.photos/100",
	})
	if err != nil {
		t.Error("New user is not created properly")
	}
}

func TestUserRepositoryStub_Count(t *testing.T) {
	repo := NewUserRepositoryStub()
	data := &models.UsersQueryRequestParams{}
	count, err := repo.Count(data)
	if err != nil {
		t.Error("Error in count method")
	}
	if count != 0 {
		t.Error("Count method is not working properly")
	}
}

func TestUserRepositoryStub_Query(t *testing.T) {
	repo := NewUserRepositoryStub()
	data := &models.UsersQueryRequestParams{}
	_, err := repo.Query(data)
	if err != nil {
		t.Error("Error in query method")
	}
}

func TestUserRepositoryStub_Update(t *testing.T) {
	repo := NewUserRepositoryStub()
	data := &models.UsersUpdateRequestBody{
		FirstName:    "TestName",
		LastName:     "TestLastName",
		Username:     "test",
		Email:        "value@test.com",
		Phone:        "+358442362057",
		NationalCode: "111111111",
		BirthDate:    "1111111",
		AvatarUrl:    "https://picsum.photos/100",
	}
	_, err := repo.Update(data, 1)
	if err != nil {
		t.Error("Error in update method")
	}
}

func TestUserRepositoryDB_Delete(t *testing.T) {
	repo := NewUserRepositoryStub()
	data := &models.UsersDeleteRequestBody{
		IDs: `[1,2,3]`,
	}
	_, err := repo.Delete(data)
	if err != nil {
		t.Error("Error in delete method")
	}
}
