package response

import "net/http"

type HandlerFunc = func(w http.ResponseWriter, r *http.Request) Response

type Handler struct {
	Fn HandlerFunc
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := h.Fn(w, r)
	switch result.(type) {
	case ErrorResponse:
		err := result.(ErrorResponse)
		JsonWithWriter(w, map[string]interface{}{
			"errors":  err.Data,
			"message": err.Message,
		}, err.StatusCode)
	case SuccessResponse:
		data := result.(SuccessResponse)
		JsonWithWriter(w, data, data.StatusCode)
	}
}

func Handle(fn HandlerFunc) Handler {
	return Handler{Fn: fn}
}
