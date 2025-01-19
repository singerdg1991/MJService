package handlers

import (
	"testing"

	languageSkillRepo "github.com/hoitek/Maja-Service/internal/languageskill/repositories"
	languageSkillService "github.com/hoitek/Maja-Service/internal/languageskill/service"
	roleRepo "github.com/hoitek/Maja-Service/internal/role/repositories"
	roleService "github.com/hoitek/Maja-Service/internal/role/service"
	"github.com/hoitek/Maja-Service/internal/user/repositories"
	"github.com/hoitek/Maja-Service/internal/user/service"
)

func TestNewUserHandler(t *testing.T) {
	t.Run("Check user handler", func(t *testing.T) {
		userRepository := repositories.NewUserRepositoryStub()
		roleRepository := roleRepo.NewRoleRepositoryStub()
		languageSkillRepository := languageSkillRepo.NewLanguageSkillRepositoryStub()
		userService := service.NewUserService(userRepository, nil, nil)
		roleService := roleService.NewRoleService(roleRepository, nil)
		languageSkillService := languageSkillService.NewLanguageSkillService(languageSkillRepository, nil)
		_, err := NewUserHandler(nil, userService, roleService, &languageSkillService, nil, nil)

		if err == nil {
			t.Error("Router can not be nil")
		}
	})
}

func TestCreate(t *testing.T) {
	//t.Run("Check error response", func(t *testing.T) {
	//	r := router.Init()
	//	userRepository := repositories.NewUserRepositoryStub()
	//	userService := service.NewUserService(userRepository, nil, nil)
	//	userHandler, err := NewUserHandler(r, userService)
	//	emptyUserHandler := UserHandler{}
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	if userHandler == emptyUserHandler {
	//		t.Error("UserHandler can not be empty")
	//	}
	//	repositories.UserTestCondition.HasError = true
	//	handler := userHandler.Create()
	//	req := httptest.NewRequest(http.MethodGet, "/", nil)
	//	res := httptest.NewRecorder()
	//	result := handler.Fn(res, req)
	//	switch result.(type) {
	//	case response.SuccessResponse:
	//		t.Error("Response can not be SuccessResponse")
	//	}
	//})
	//
	//t.Run("Check error response when body is invalid", func(t *testing.T) {
	//	r := router.Init()
	//	userRepository := repositories.NewUserRepositoryStub()
	//	userService := service.NewUserService(userRepository, nil, nil)
	//	userHandler, err := NewUserHandler(r, userService)
	//	emptyUserHandler := UserHandler{}
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	if userHandler == emptyUserHandler {
	//		t.Error("UserHandler can not be empty")
	//	}
	//	handler := userHandler.Create()
	//	body := strings.NewReader(`{"wrongKey":"value"}`)
	//	req := httptest.NewRequest(http.MethodGet, "/", body)
	//	res := httptest.NewRecorder()
	//	result := handler.Fn(res, req)
	//	switch result.(type) {
	//	case response.SuccessResponse:
	//		t.Error("Response can not be SuccessResponse")
	//	}
	//})

	//t.Run("Check success response", func(t *testing.T) {
	//	r := router.Init()
	//	rDB := repositories.NewUserRepositoryStub()
	//	rGRPC := repositories.NewUserRepositoryGRPC(grpc.Connection)
	//	userService := service.NewUserService(rDB, rGRPC)
	//	userHandler, err := NewUserHandler(r, userService)
	//	emptyUserHandler := UserHandler{}
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	if userHandler == emptyUserHandler {
	//		t.Error("UserHandler can not be empty")
	//	}
	//	handler := userHandler.Create()
	//	body := strings.NewReader(`{"name":"value","lastName":"value","username":"value","email":"value@domain.com","phone":"value","nationalCode":"value","birthDate":"value","avatarUrl":"https://value.com"}`)
	//	req := httptest.NewRequest(http.MethodGet, "/", body)
	//	res := httptest.NewRecorder()
	//	result := handler.Fn(res, req)
	//	log.Println(result)
	//	switch result.(type) {
	//	case response.ErrorResponse:
	//		t.Error("Response can not be ErrorResponse")
	//	}
	//})
}
