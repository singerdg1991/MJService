package router

import (
	"net/http"

	"github.com/hoitek/Maja-Service/config"

	"github.com/gorilla/mux"
	bodylimit "github.com/hoitek/Middlewares/bodyLimit"
	"github.com/hoitek/Middlewares/gzip"
)

var Default *mux.Router

func Init() *mux.Router {
	Default = mux.NewRouter()
	Default.Use(bodylimit.Middleware(config.AppConfig.MaxBodySizeLimit))
	Default.Use(gzip.Middleware(1))
	Default.Use(CORSMiddleware)
	return Default
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}
