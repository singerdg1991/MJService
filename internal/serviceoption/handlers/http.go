package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/serviceoption/models"
	"github.com/hoitek/Maja-Service/internal/serviceoption/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/serviceoption/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type ServiceOptionHandler struct {
	UserService          uPorts.UserService
	ServiceOptionService ports.ServiceOptionService
}

func NewServiceOptionHandler(r *mux.Router, s ports.ServiceOptionService, u uPorts.UserService) (ServiceOptionHandler, error) {
	serviceOptionHandler := ServiceOptionHandler{
		UserService:          u,
		ServiceOptionService: s,
	}
	if r == nil {
		return ServiceOptionHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ServiceOptionConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ServiceOptionConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/serviceoptions",
			Method:      http.MethodPost,
			Handler:     serviceOptionHandler.Create(),
			Permissions: []string{permissions.SERVICE_OPTIONS_CREATE_NEW_SERVICE_OPTION},
		},
		{
			Path:        "/serviceoptions",
			Method:      http.MethodGet,
			Handler:     serviceOptionHandler.Query(),
			Permissions: []string{permissions.SERVICE_OPTIONS_VIEW_ALL_SERVICE_OPTIONS},
		},
		{
			Path:        "/serviceoptions",
			Method:      http.MethodDelete,
			Handler:     serviceOptionHandler.Delete(),
			Permissions: []string{permissions.SERVICE_OPTIONS_CREATE_NEW_SERVICE_OPTION},
		},
		{
			Path:        "/serviceoptions/{id}",
			Method:      http.MethodPut,
			Handler:     serviceOptionHandler.Update(),
			Permissions: []string{permissions.SERVICE_OPTIONS_CREATE_NEW_SERVICE_OPTION},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return serviceOptionHandler, nil
}

/*
* @apiTag: serviceoption
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: serviceoptions
* @apiResponseRef: ServiceOptionsQueryResponse
* @apiSummary: Query Services
* @apiParametersRef: ServiceOptionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ServiceOptionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ServiceOptionsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ServiceOptionHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ServiceOptionsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query serviceoptions
		serviceoptions, err := h.ServiceOptionService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(serviceoptions)
	})
}

/*
 * @apiTag: serviceoption
 * @apiPath: /serviceoptions
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: ServiceOptionsCreateRequestBody
 * @apiResponseRef: ServiceOptionsCreateResponse
 * @apiSummary: Create serviceOption
 * @apiDescription: Create serviceOption
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceOptionHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServiceOptionsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find service option by name and service type id
		res, _ := h.ServiceOptionService.FindByNameAndServiceTypeID(payload.Name, int(payload.ServiceTypeID))
		if res != nil {
			return response.ErrorInternalServerError(nil, "Service option already exists")
		}

		// Create service option
		serviceOption, err := h.ServiceOptionService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(serviceOption)
	})
}

/*
 * @apiTag: serviceoption
 * @apiPath: /serviceoptions/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: ServiceOptionsUpdateRequestParams
 * @apiRequestRef: ServiceOptionsCreateRequestBody
 * @apiResponseRef: ServiceOptionsCreateResponse
 * @apiSummary: Update serviceOption
 * @apiDescription: Update serviceOption
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceOptionHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.ServiceOptionsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.ServiceOptionsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Service already exists
		res, _ := h.ServiceOptionService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Service does not exist")
		}

		// Check if Service already exists with the same name
		res, _ = h.ServiceOptionService.FindByNameAndServiceTypeID(payload.Name, int(payload.ServiceTypeID))
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Service already exists")
		}

		// Update Service
		data, err := h.ServiceOptionService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: serviceoption
 * @apiPath: /serviceoptions
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: ServiceOptionsDeleteRequestBody
 * @apiResponseRef: ServiceOptionsDeleteResponse
 * @apiSummary: Delete serviceOption
 * @apiDescription: Delete serviceOption
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceOptionHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServiceOptionsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.ServiceOptionService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
