package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/license/models"
	"github.com/hoitek/Maja-Service/internal/license/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/license/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type LicenseHandler struct {
	UserService    uPorts.UserService
	LicenseService ports.LicenseService
}

func NewLicenseHandler(r *mux.Router, s ports.LicenseService, u uPorts.UserService) (LicenseHandler, error) {
	licenseHandler := LicenseHandler{
		UserService:    u,
		LicenseService: s,
	}
	if r == nil {
		return LicenseHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.LicenseConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.LicenseConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/licenses",
			Method:      http.MethodPost,
			Handler:     licenseHandler.Create(),
			Permissions: []string{permissions.LICENSES_CREATE_NEW_LICENSE},
		},
		{
			Path:        "/licenses",
			Method:      http.MethodGet,
			Handler:     licenseHandler.Query(),
			Permissions: []string{permissions.LICENSES_VIEW_ALL_LICENSES},
		},
		{
			Path:        "/licenses",
			Method:      http.MethodDelete,
			Handler:     licenseHandler.Delete(),
			Permissions: []string{permissions.LICENSES_CREATE_NEW_LICENSE},
		},
		{
			Path:        "/licenses/{id}",
			Method:      http.MethodPut,
			Handler:     licenseHandler.Update(),
			Permissions: []string{permissions.LICENSES_CREATE_NEW_LICENSE},
		},
		{
			Path:        "/licenses/csv/download",
			Method:      http.MethodGet,
			Handler:     licenseHandler.Download(),
			Permissions: []string{permissions.LICENSES_VIEW_ALL_LICENSES},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return licenseHandler, nil
}

/*
* @apiTag: license
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: licenses
* @apiResponseRef: LicensesQueryResponse
* @apiSummary: Query licenses
* @apiParametersRef: LicensesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: LicensesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: LicensesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *LicenseHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.LicensesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query licenses
		licenses, err := h.LicenseService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(licenses)
	})
}

/*
* @apiTag: license
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: licenses/csv/download
* @apiResponseRef: LicensesQueryResponse
* @apiSummary: Query licenses
* @apiParametersRef: LicensesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: LicensesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: LicensesQueryNotFoundResponse
 */
func (h *LicenseHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.LicensesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		licenses, err := h.LicenseService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(licenses)
	})
}

/*
 * @apiTag: license
 * @apiPath: /licenses
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: LicensesCreateRequestBody
 * @apiResponseRef: LicensesCreateResponse
 * @apiSummary: Create license
 * @apiDescription: Create license
 * @apiSecurity: apiKeySecurity
 */
func (h *LicenseHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.LicensesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find license by name
		res, _ := h.LicenseService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "License already exists")
		}

		// Create license
		license, err := h.LicenseService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(license)
	})
}

/*
 * @apiTag: license
 * @apiPath: /licenses/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: LicensesUpdateRequestParams
 * @apiRequestRef: LicensesCreateRequestBody
 * @apiResponseRef: LicensesCreateResponse
 * @apiSummary: Update license
 * @apiDescription: Update license
 * @apiSecurity: apiKeySecurity
 */
func (h *LicenseHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.LicensesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.LicensesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.LicenseService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.LicenseService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.LicenseService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: license
 * @apiPath: /licenses
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: LicensesDeleteRequestBody
 * @apiResponseRef: LicensesDeleteResponse
 * @apiSummary: Delete license
 * @apiDescription: Delete license
 * @apiSecurity: apiKeySecurity
 */
func (h *LicenseHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.LicensesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.LicenseService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
