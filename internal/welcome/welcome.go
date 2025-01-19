package welcome

import (
	"github.com/gorilla/mux"
	"github.com/hoitek/Maja-Service/internal/welcome/config"
	"github.com/hoitek/Maja-Service/internal/welcome/handlers"
)

type module struct {
	Config config.ConfigType
}

var Module = &module{}

func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	config.WelcomeConfig = &c
	return m
}

func (m *module) RegisterHTTP(r *mux.Router) (handlers.WelcomeHandler, error) {
	r = r.PathPrefix("/").Subrouter()
	handler, err := handlers.NewWelcomeHandler(r)
	if err != nil {
		return handlers.WelcomeHandler{}, err
	}
	return handler, nil
}
