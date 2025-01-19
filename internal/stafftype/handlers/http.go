package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/stafftype/models"
	"github.com/hoitek/Maja-Service/internal/stafftype/ports"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/stafftype/config"
)

type StaffTypeHandler struct {
	StaffTypeService ports.StaffTypeService
}

func NewStaffTypeHandler(r *mux.Router, s ports.StaffTypeService) (StaffTypeHandler, error) {
	staffTypeHandler := StaffTypeHandler{
		StaffTypeService: s,
	}
	if r == nil {
		return StaffTypeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.StaffTypeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.StaffTypeConfig.ApiVersion1).Subrouter()

	rv1.Handle("/stafftypes", staffTypeHandler.Query()).Methods(http.MethodGet)

	return staffTypeHandler, nil
}

/*
* @apiTag: stafftype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: stafftypes
* @apiResponseRef: StaffTypesQueryResponse
* @apiSummary: Query staff types
* @apiParametersRef: StaffTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: StaffTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: StaffTypesQueryNotFoundResponse
 */
func (h *StaffTypeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.StaffTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		staffTypes, err := h.StaffTypeService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(staffTypes)
	})
}
