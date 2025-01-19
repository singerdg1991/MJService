package handlers

import (
	"errors"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/staffclub/warning/domain"

	"github.com/hoitek/Maja-Service/internal/_shared/utils"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	rPorts "github.com/hoitek/Maja-Service/internal/punishment/ports"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/models"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type WarningHandler struct {
	UserService       uPorts.UserService
	WarningService    ports.WarningService
	PunishmentService rPorts.PunishmentService
}

func NewWarningHandler(r *mux.Router, s ports.WarningService, rs rPorts.PunishmentService, us uPorts.UserService) (WarningHandler, error) {
	warningHandler := WarningHandler{
		UserService:       us,
		WarningService:    s,
		PunishmentService: rs,
	}
	if r == nil {
		return WarningHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.WarningConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.WarningConfig.ApiVersion1).Subrouter()

	// Add JWT middleware
	rAuth := rv1.PathPrefix("/").Subrouter()
	rAuth.Use(middlewares.OAuth2Middleware)
	rAuth.Use(middlewares.AuthMiddleware(us, []string{}))

	rAuth.Handle("/staffclub/warnings", warningHandler.Create()).Methods(http.MethodPost)
	rAuth.Handle("/staffclub/warnings", warningHandler.Query()).Methods(http.MethodGet)
	rAuth.Handle("/staffclub/warnings", warningHandler.Delete()).Methods(http.MethodDelete)
	rAuth.Handle("/staffclub/warnings/{id}", warningHandler.Update()).Methods(http.MethodPut)
	rAuth.Handle("/staffclub/warnings/csv/download", warningHandler.Download()).Methods(http.MethodGet)

	return warningHandler, nil
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/warnings
* @apiResponseRef: WarningsQueryResponse
* @apiSummary: Query warnings
* @apiParametersRef: WarningsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: WarningsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: WarningsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *WarningHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.WarningsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query warnings
		warnings, err := h.WarningService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(warnings)
	})
}

/*
* @apiTag: staffclub
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: /staffclub/warnings/csv/download
* @apiResponseRef: WarningsQueryResponse
* @apiSummary: Query warnings
* @apiParametersRef: WarningsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: WarningsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: WarningsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *WarningHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.WarningsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		warnings, err := h.WarningService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(warnings)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/warnings
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: WarningsCreateRequestBody
 * @apiResponseRef: WarningsCreateResponse
 * @apiSummary: Create warning
 * @apiDescription: Create warning
 * @apiSecurity: apiKeySecurity
 */
func (h *WarningHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.WarningsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Warning already exists with the same warning number
		warningByNumber, _ := h.WarningService.FindByWarningNumber(int64(payload.WarningNumber))
		if warningByNumber != nil {
			return response.ErrorBadRequest(nil, "Warning with this warning number already exists")
		}

		// Find punishment
		punishment, _ := h.PunishmentService.FindByID(int64(payload.PunishmentID))
		if punishment == nil {
			return response.ErrorInternalServerError(nil, "Punishment not found")
		}
		payload.Punishment = &domain.WarningPunishment{
			ID:          punishment.ID,
			Name:        punishment.Name,
			Description: *punishment.Description,
		}

		// Create warning
		warning, err := h.WarningService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(warning)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/warnings/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: WarningsUpdateRequestParams
 * @apiRequestRef: WarningsCreateRequestBody
 * @apiResponseRef: WarningsCreateResponse
 * @apiSummary: Update warning
 * @apiDescription: Update warning
 * @apiSecurity: apiKeySecurity
 */
func (h *WarningHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.WarningsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.WarningsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Warning already exists
		warning, _ := h.WarningService.FindByID(int64(params.ID))
		if warning == nil {
			return response.ErrorBadRequest(nil, "Warning does not exist")
		}

		// Check if Warning already exists with the same warning number
		warningByNumber, _ := h.WarningService.FindByWarningNumber(int64(warning.WarningNumber))
		if warningByNumber != nil && int64(warningByNumber.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Warning with this warning number already exists")
		}

		// Find punishment
		punishment, _ := h.PunishmentService.FindByID(int64(payload.PunishmentID))
		if punishment == nil {
			return response.ErrorInternalServerError(nil, "Punishment not found")
		}
		payload.Punishment = &domain.WarningPunishment{
			ID:          punishment.ID,
			Name:        punishment.Name,
			Description: *punishment.Description,
		}

		// Update Warning
		data, err := h.WarningService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: staffclub
 * @apiPath: /staffclub/warnings
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: WarningsDeleteRequestBody
 * @apiResponseRef: WarningsDeleteResponse
 * @apiSummary: Delete warning
 * @apiDescription: Delete warning
 * @apiSecurity: apiKeySecurity
 */
func (h *WarningHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.WarningsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.WarningService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
