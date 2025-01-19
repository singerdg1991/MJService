package response

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	ErrorResponse | SuccessResponse | any
}

func JsonWithWriter[T Response](w http.ResponseWriter, data T, statusCode int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func SuccessWithWriter(w http.ResponseWriter, data SuccessResponse) error {
	return JsonWithWriter(w, data, http.StatusOK)
}

func ErrorWithWriter(w http.ResponseWriter, errorResponse ErrorResponse, statusCode int, messages ...string) error {
	if len(messages) > 0 {
		errorResponse.Message = messages[0]
	}
	return JsonWithWriter(w, errorResponse, statusCode)
}
