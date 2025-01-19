package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type WelcomeHandler struct{}

func NewWelcomeHandler(r *mux.Router) (WelcomeHandler, error) {
	welcomeHandler := WelcomeHandler{}
	if r == nil {
		return WelcomeHandler{}, errors.New("router can not be nil")
	}
	r.HandleFunc("/", welcomeHandler.WelcomeHome).Methods(http.MethodGet)
	r.HandleFunc("/api", welcomeHandler.WelcomeApi).Methods(http.MethodGet)
	r.HandleFunc("/api/v1", welcomeHandler.WelcomeApiV1).Methods(http.MethodGet)
	r.HandleFunc("/api/v2", welcomeHandler.WelcomeApiV2).Methods(http.MethodGet)
	return welcomeHandler, nil
}

func (h *WelcomeHandler) WelcomeHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Maja!"))
}

func (h *WelcomeHandler) WelcomeApi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Maja Api Server!"))
}

func (h *WelcomeHandler) WelcomeApiV1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Maja Api V1 Server!"))
}

func (h *WelcomeHandler) WelcomeApiV2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Maja Api V2 Server!"))
}
