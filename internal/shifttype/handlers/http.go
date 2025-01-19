package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/shifttype/models"
	"github.com/hoitek/Maja-Service/internal/shifttype/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/shifttype/config"
)

type ShiftTypeHandler struct {
	ShiftTypeService service.ShiftTypeService
}

func NewShiftTypeHandler(r *mux.Router, s service.ShiftTypeService) (ShiftTypeHandler, error) {
	shiftTypeHandler := ShiftTypeHandler{
		ShiftTypeService: s,
	}
	if r == nil {
		return ShiftTypeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ShiftTypeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ShiftTypeConfig.ApiVersion1).Subrouter()

	rv1.Handle("/shift-types", shiftTypeHandler.Query()).Methods(http.MethodGet)
	return shiftTypeHandler, nil
}

/*
* @apiTag: shifttype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /shift-types
* @apiResponseRef: ShiftTypesQueryResponse
* @apiSummary: Query ShiftTypes
* @apiParametersRef: ShiftTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ShiftTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ShiftTypesQueryNotFoundResponse
 */
func (h *ShiftTypeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ShiftTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		shiftTypes, err := h.ShiftTypeService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(shiftTypes)
	})
}
