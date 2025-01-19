package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/limitation/models"
	"github.com/hoitek/Maja-Service/internal/limitation/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/limitation/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type LimitationHandler struct {
	UserService       uPorts.UserService
	LimitationService ports.LimitationService
}

func NewLimitationHandler(r *mux.Router, s ports.LimitationService, u uPorts.UserService) (LimitationHandler, error) {
	limitationHandler := LimitationHandler{
		UserService:       u,
		LimitationService: s,
	}
	if r == nil {
		return LimitationHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.LimitationConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.LimitationConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/limitations",
			Method:      http.MethodPost,
			Handler:     limitationHandler.Create(),
			Permissions: []string{permissions.LIMITATIONS_CREATE_NEW_LIMITATION},
		},
		{
			Path:        "/limitations",
			Method:      http.MethodGet,
			Handler:     limitationHandler.Query(),
			Permissions: []string{permissions.LIMITATIONS_VIEW_ALL_LIMITATIONS},
		},
		{
			Path:        "/limitations",
			Method:      http.MethodDelete,
			Handler:     limitationHandler.Delete(),
			Permissions: []string{permissions.LIMITATIONS_CREATE_NEW_LIMITATION},
		},
		{
			Path:        "/limitations/{id}",
			Method:      http.MethodPut,
			Handler:     limitationHandler.Update(),
			Permissions: []string{permissions.LIMITATIONS_CREATE_NEW_LIMITATION},
		},
		{
			Path:        "/limitations/csv/download",
			Method:      http.MethodGet,
			Handler:     limitationHandler.Download(),
			Permissions: []string{permissions.LIMITATIONS_VIEW_ALL_LIMITATIONS},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return limitationHandler, nil
}

/*
* @apiTag: limitation
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: limitations
* @apiResponseRef: LimitationsQueryResponse
* @apiSummary: Query limitations
* @apiParametersRef: LimitationsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: LimitationsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: LimitationsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *LimitationHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.LimitationsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query limitations
		limitations, err := h.LimitationService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(limitations)
	})
}

/*
* @apiTag: limitation
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: limitations/csv/download
* @apiResponseRef: LimitationsQueryResponse
* @apiSummary: Query limitations
* @apiParametersRef: LimitationsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: LimitationsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: LimitationsQueryNotFoundResponse
 */
func (h *LimitationHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.LimitationsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		limitations, err := h.LimitationService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(limitations)
	})
}

/*
 * @apiTag: limitation
 * @apiPath: /limitations
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: LimitationsCreateRequestBody
 * @apiResponseRef: LimitationsCreateResponse
 * @apiSummary: Create limitation
 * @apiDescription: Create limitation
 * @apiSecurity: apiKeySecurity
 */
func (h *LimitationHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.LimitationsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find limitation by name
		res, _ := h.LimitationService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "Limitation already exists")
		}

		// Create limitation
		limitation, err := h.LimitationService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(limitation)
	})
}

/*
 * @apiTag: limitation
 * @apiPath: /limitations/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: LimitationsUpdateRequestParams
 * @apiRequestRef: LimitationsCreateRequestBody
 * @apiResponseRef: LimitationsCreateResponse
 * @apiSummary: Update limitation
 * @apiDescription: Update limitation
 * @apiSecurity: apiKeySecurity
 */
func (h *LimitationHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.LimitationsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.LimitationsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if limitation already exists
		res, _ := h.LimitationService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "limitation does not exist")
		}

		// Check if limitation already exists with the same name
		res, _ = h.LimitationService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "limitation already exists")
		}

		// Update limitation
		data, err := h.LimitationService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: limitation
 * @apiPath: /limitations
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: LimitationsDeleteRequestBody
 * @apiResponseRef: LimitationsDeleteResponse
 * @apiSummary: Delete limitation
 * @apiDescription: Delete limitation
 * @apiSecurity: apiKeySecurity
 */
func (h *LimitationHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.LimitationsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.LimitationService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
