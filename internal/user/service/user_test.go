package service

import (
	"context"
	"testing"

	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/repositories"
)

func TestUserServiceNewUserService(t *testing.T) {
	t.Run("Check service when repository is ready", func(t *testing.T) {
		rDB := repositories.NewUserRepositoryStub()
		s := NewUserService(rDB, nil, nil)
		if s == nil {
			t.Error("Repository can not be nil")
		}
	})
}

func TestUserServiceBindUserToContext(t *testing.T) {
	t.Run("Check context result", func(t *testing.T) {
		rDB := repositories.NewUserRepositoryStub()
		s := NewUserService(rDB, nil, nil)
		ctx := s.BindUserToContext(context.Background(), &domain.User{})
		if ctx == nil {
			t.Error("Context can not be nil")
		}
	})

	t.Run("When context is nil", func(t *testing.T) {
		rDB := repositories.NewUserRepositoryStub()
		s := NewUserService(rDB, nil, nil)
		ctx := s.BindUserToContext(context.TODO(), &domain.User{})
		if ctx == nil {
			t.Error("Context can not be nil")
		}
	})
}

func TestCreate(t *testing.T) {
	//rDB := repositories.NewUserRepositoryStub()
	//rGRPC := repositories.NewUserRepositoryGRPC(grpc.Connection)
	//s := NewUserService(rDB, rGRPC)
	//_, err := s.Create(&models.UsersCreateRequestBody{
	//	Name:         "TestName",
	//	LastName:     "TestLastName",
	//	Username:     "test",
	//	Email:        "test@yahoo.com",
	//	Phone:        "+358442362057",
	//	NationalCode: "111111111",
	//	BirthDate:    "111111111",
	//	AvatarUrl:    "https://picsum.photos/100",
	//})
	//
	//if err != nil {
	//	t.Error("Can't return error")
	//}
}
