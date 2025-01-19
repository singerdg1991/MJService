package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	s3Ports "github.com/hoitek/Maja-Service/internal/s3/ports"
	"github.com/hoitek/Maja-Service/internal/static/config"
	"github.com/hoitek/Maja-Service/internal/static/ports"
	"io"
	"log"
	"net/http"
)

type StaticHandler struct {
	StaticService ports.StaticService
	S3Service     s3Ports.S3Service
}

func NewStaticHandler(r *mux.Router, ss ports.StaticService, s3 s3Ports.S3Service) (StaticHandler, error) {
	staticHandler := StaticHandler{
		StaticService: ss,
		S3Service:     s3,
	}
	if r == nil {
		return StaticHandler{}, errors.New("router can not be nil")
	}
	r.HandleFunc("/uploads/{domain}/{file}", staticHandler.FileHandler)
	r.PathPrefix("/").Handler(http.StripPrefix(config.StaticConfig.StaticURIPath, http.FileServer(http.Dir(config.StaticConfig.StaticDir))))
	return staticHandler, nil
}

func (h *StaticHandler) FileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["file"]
	bucketName := fmt.Sprintf("%s.%s", "maja", vars["domain"])
	cacheKey := h.S3Service.GetCacheKey(bucketName, fileName)

	// Check cache
	data, err := h.StaticService.GetCache(cacheKey)
	if err == nil {
		w.Write(data)
		return
	}

	// Get object
	obj, err := h.StaticService.GetObject(bucketName, fileName)
	if err != nil {
		log.Println(err)
		http.Error(w, errors.New("file not found").Error(), http.StatusNotFound)
		return
	}

	// Read the file data
	fileBytes, err := io.ReadAll(obj)
	if err != nil {
		log.Println(err)
		http.Error(w, errors.New("file not found").Error(), http.StatusNotFound)
		return
	}

	// Set cache
	err = h.StaticService.SetCache(cacheKey, fileBytes)
	if err != nil {
		log.Println(err)
		http.Error(w, errors.New("file not found").Error(), http.StatusNotFound)
		return
	}

	// Send object to response writer if object is exist
	w.Write(fileBytes)
}
