package user

import (
	"github.com/hoitek/Maja-Service/database"
	"testing"

	"github.com/hoitek/Maja-Service/internal/user/config"
	"github.com/hoitek/Maja-Service/router"
)

func TestUserGetService(t *testing.T) {
	t.Run("When repository is nil", func(t *testing.T) {
		var m *module = &module{}
		s := m.GetUserService(database.PostgresDB)
		if s == nil {
			t.Error("Service does not create properly")
		}
	})
}

func TestUserSetConfig(t *testing.T) {
	var m *module = &module{}
	m.Setup(config.ConfigType{
		Environment: "this is test",
	}).SetDatabases(database.PostgresDB, database.MongoDB)
	if m.Config.Environment != "this is test" {
		t.Error("Config does not set")
	}
}

func TestUserRegisterHTTP(t *testing.T) {
	t.Run("When router is nil", func(t *testing.T) {
		var m *module = &module{}
		_, err := m.RegisterHTTP(nil)
		if err == nil {
			t.Error("Router can not be nil")
		}
	})

	t.Run("When router is not nil", func(t *testing.T) {
		var m *module = &module{}
		r := router.Init()
		_, err := m.RegisterHTTP(r)
		if err != nil {
			t.Error("Error should be nil")
		}
	})
}
