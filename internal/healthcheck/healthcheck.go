package healthcheck

import (
	"github.com/gorilla/mux"
	"github.com/hoitek/Maja-Service/internal/healthcheck/config"
	"github.com/hoitek/Maja-Service/internal/healthcheck/handlers"
)

type module struct {
	Config config.ConfigType
}

var Module = &module{}

func (m *module) Setup(c config.ConfigType) *module {
	m.Config = c
	return m
}

func (m *module) RegisterHTTP(r *mux.Router) (handlers.HealthCheckHandler, error) {
	r = r.PathPrefix("/").Subrouter()
	handler, err := handlers.NewHealthCheckHandler(r)
	if err != nil {
		return handlers.HealthCheckHandler{}, err
	}
	return handler, nil
}
