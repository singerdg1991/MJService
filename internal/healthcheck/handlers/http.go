package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type HealthCheckHandler struct{}

func NewHealthCheckHandler(r *mux.Router) (HealthCheckHandler, error) {
	healthCheckHandler := HealthCheckHandler{}
	if r == nil {
		return HealthCheckHandler{}, errors.New("router can not be nil")
	}
	r.HandleFunc("/healthcheck", healthCheckHandler.Home).Methods(http.MethodGet)
	return healthCheckHandler, nil
}

func (h *HealthCheckHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
