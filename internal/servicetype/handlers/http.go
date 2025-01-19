package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/servicetype/models"
	"github.com/hoitek/Maja-Service/internal/servicetype/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/servicetype/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type ServiceTypeHandler struct {
	UserService        uPorts.UserService
	ServiceTypeService ports.ServiceTypeService
}

func NewServiceTypeHandler(r *mux.Router, s ports.ServiceTypeService, u uPorts.UserService) (ServiceTypeHandler, error) {
	serviceTypeHandler := ServiceTypeHandler{
		UserService:        u,
		ServiceTypeService: s,
	}
	if r == nil {
		return ServiceTypeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ServiceTypeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ServiceTypeConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/servicetypes",
			Method:      http.MethodPost,
			Handler:     serviceTypeHandler.Create(),
			Permissions: []string{permissions.SERVICE_TYPES_CREATE_NEW_SERVICE_TYPE},
		},
		{
			Path:        "/servicetypes",
			Method:      http.MethodGet,
			Handler:     serviceTypeHandler.Query(),
			Permissions: []string{permissions.SERVICE_TYPES_VIEW_ALL_SERVICE_TYPES},
		},
		{
			Path:        "/servicetypes",
			Method:      http.MethodDelete,
			Handler:     serviceTypeHandler.Delete(),
			Permissions: []string{permissions.SERVICE_TYPES_CREATE_NEW_SERVICE_TYPE},
		},
		{
			Path:        "/servicetypes/{id}",
			Method:      http.MethodPut,
			Handler:     serviceTypeHandler.Update(),
			Permissions: []string{permissions.SERVICE_TYPES_CREATE_NEW_SERVICE_TYPE},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return serviceTypeHandler, nil
}

/*
* @apiTag: servicetype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: servicetypes
* @apiResponseRef: ServiceTypesQueryResponse
* @apiSummary: Query Report Categories
* @apiParametersRef: ServiceTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ServiceTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ServiceTypesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ServiceTypeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ServiceTypesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query report categories
		serviceTypes, err := h.ServiceTypeService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(serviceTypes)
	})
}

/*
 * @apiTag: servicetype
 * @apiPath: /servicetypes
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: ServiceTypesCreateRequestBody
 * @apiResponseRef: ServiceTypesCreateResponse
 * @apiSummary: Create serviceType
 * @apiDescription: Create serviceType
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceTypeHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServiceTypesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find report category by name and serviceID
		res, _ := h.ServiceTypeService.FindByNameAndServiceID(payload.Name, int(payload.ServiceID))
		if res != nil {
			return response.ErrorInternalServerError(nil, "ServiceType already exists")
		}

		// Create report category
		serviceType, err := h.ServiceTypeService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(serviceType)
	})
}

/*
 * @apiTag: servicetype
 * @apiPath: /servicetypes/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: ServiceTypesUpdateRequestParams
 * @apiRequestRef: ServiceTypesCreateRequestBody
 * @apiResponseRef: ServiceTypesCreateResponse
 * @apiSummary: Update serviceType
 * @apiDescription: Update serviceType
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceTypeHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.ServiceTypesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.ServiceTypesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Report Category already exists
		res, _ := h.ServiceTypeService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Report Category does not exist")
		}

		// Check if Report Category already exists with the same name
		res, _ = h.ServiceTypeService.FindByNameAndServiceID(payload.Name, int(payload.ServiceID))
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Report Category already exists")
		}

		// Update Report Category
		data, err := h.ServiceTypeService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: servicetype
 * @apiPath: /servicetypes
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: ServiceTypesDeleteRequestBody
 * @apiResponseRef: ServiceTypesDeleteResponse
 * @apiSummary: Delete serviceType
 * @apiDescription: Delete serviceType
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceTypeHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServiceTypesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.ServiceTypeService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
