package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/service/models"
	"github.com/hoitek/Maja-Service/internal/service/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/service/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type ServiceHandler struct {
	UserService    uPorts.UserService
	ServiceService ports.ServiceService
}

func NewServiceHandler(r *mux.Router, s ports.ServiceService, u uPorts.UserService) (ServiceHandler, error) {
	serviceHandler := ServiceHandler{
		UserService:    u,
		ServiceService: s,
	}
	if r == nil {
		return ServiceHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ServiceConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ServiceConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/services",
			Method:      http.MethodPost,
			Handler:     serviceHandler.Create(),
			Permissions: []string{permissions.SERVICES_CREATE_NEW_SERVICE},
		},
		{
			Path:        "/services",
			Method:      http.MethodGet,
			Handler:     serviceHandler.Query(),
			Permissions: []string{permissions.SERVICES_VIEW_ALL_SERVICES},
		},
		{
			Path:        "/services",
			Method:      http.MethodDelete,
			Handler:     serviceHandler.Delete(),
			Permissions: []string{permissions.SERVICES_CREATE_NEW_SERVICE},
		},
		{
			Path:        "/services/{id}",
			Method:      http.MethodPut,
			Handler:     serviceHandler.Update(),
			Permissions: []string{permissions.SERVICES_CREATE_NEW_SERVICE},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return serviceHandler, nil
}

/*
* @apiTag: service
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: services
* @apiResponseRef: ServicesQueryResponse
* @apiSummary: Query Report Types
* @apiParametersRef: ServicesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ServicesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ServicesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ServiceHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ServicesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query report types
		services, err := h.ServiceService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(services)
	})
}

/*
 * @apiTag: service
 * @apiPath: /services
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: ServicesCreateRequestBody
 * @apiResponseRef: ServicesCreateResponse
 * @apiSummary: Create service
 * @apiDescription: Create service
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServicesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find service by name
		res, _ := h.ServiceService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "Service already exists")
		}

		// Create service
		service, err := h.ServiceService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(service)
	})
}

/*
 * @apiTag: service
 * @apiPath: /services/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: ServicesUpdateRequestParams
 * @apiRequestRef: ServicesCreateRequestBody
 * @apiResponseRef: ServicesCreateResponse
 * @apiSummary: Update service
 * @apiDescription: Update service
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.ServicesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.ServicesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Report Type already exists
		res, _ := h.ServiceService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Report Type does not exist")
		}

		// Check if Report Type already exists with the same name
		res, _ = h.ServiceService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Report Type already exists")
		}

		// Update Report Type
		data, err := h.ServiceService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: service
 * @apiPath: /services
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: ServicesDeleteRequestBody
 * @apiResponseRef: ServicesDeleteResponse
 * @apiSummary: Delete service
 * @apiDescription: Delete service
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServicesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.ServiceService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
