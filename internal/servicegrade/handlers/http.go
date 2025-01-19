package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/servicegrade/models"
	"github.com/hoitek/Maja-Service/internal/servicegrade/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/servicegrade/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type ServiceGradeHandler struct {
	UserService         uPorts.UserService
	ServiceGradeService ports.ServiceGradeService
}

func NewServiceGradeHandler(r *mux.Router, s ports.ServiceGradeService, u uPorts.UserService) (ServiceGradeHandler, error) {
	servicegradeHandler := ServiceGradeHandler{
		UserService:         u,
		ServiceGradeService: s,
	}
	if r == nil {
		return ServiceGradeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ServiceGradeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ServiceGradeConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/servicegrades",
			Method:      http.MethodPost,
			Handler:     servicegradeHandler.Create(),
			Permissions: []string{permissions.SERVICE_GRADES_CREATE_NEW_SERVICE_GRADE},
		},
		{
			Path:        "/servicegrades",
			Method:      http.MethodGet,
			Handler:     servicegradeHandler.Query(),
			Permissions: []string{permissions.SERVICE_GRADES_VIEW_ALL_SERVICE_GRADES},
		},
		{
			Path:        "/servicegrades",
			Method:      http.MethodDelete,
			Handler:     servicegradeHandler.Delete(),
			Permissions: []string{permissions.SERVICE_GRADES_CREATE_NEW_SERVICE_GRADE},
		},
		{
			Path:        "/servicegrades/{id}",
			Method:      http.MethodPut,
			Handler:     servicegradeHandler.Update(),
			Permissions: []string{permissions.SERVICE_GRADES_CREATE_NEW_SERVICE_GRADE},
		},
		{
			Path:        "/servicegrades/csv/download",
			Method:      http.MethodGet,
			Handler:     servicegradeHandler.Download(),
			Permissions: []string{permissions.SERVICE_GRADES_VIEW_ALL_SERVICE_GRADES},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return servicegradeHandler, nil
}

/*
* @apiTag: servicegrade
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: servicegrades
* @apiResponseRef: ServiceGradesQueryResponse
* @apiSummary: Query servicegrades
* @apiParametersRef: ServiceGradesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ServiceGradesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ServiceGradesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ServiceGradeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ServiceGradesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query servicegrades
		servicegrades, err := h.ServiceGradeService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(servicegrades)
	})
}

/*
* @apiTag: servicegrade
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: servicegrades/csv/download
* @apiResponseRef: ServiceGradesQueryResponse
* @apiSummary: Query servicegrades
* @apiParametersRef: ServiceGradesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ServiceGradesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ServiceGradesQueryNotFoundResponse
 */
func (h *ServiceGradeHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ServiceGradesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		servicegrades, err := h.ServiceGradeService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(servicegrades)
	})
}

/*
 * @apiTag: servicegrade
 * @apiPath: /servicegrades
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: ServiceGradesCreateRequestBody
 * @apiResponseRef: ServiceGradesCreateResponse
 * @apiSummary: Create servicegrade
 * @apiDescription: Create servicegrade
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceGradeHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServiceGradesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find servicegrade by name
		res, _ := h.ServiceGradeService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "ServiceGrade already exists")
		}

		// Create servicegrade
		servicegrade, err := h.ServiceGradeService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(servicegrade)
	})
}

/*
 * @apiTag: servicegrade
 * @apiPath: /servicegrades/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: ServiceGradesUpdateRequestParams
 * @apiRequestRef: ServiceGradesCreateRequestBody
 * @apiResponseRef: ServiceGradesCreateResponse
 * @apiSummary: Update servicegrade
 * @apiDescription: Update servicegrade
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceGradeHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.ServiceGradesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.ServiceGradesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.ServiceGradeService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.ServiceGradeService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.ServiceGradeService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: servicegrade
 * @apiPath: /servicegrades
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: ServiceGradesDeleteRequestBody
 * @apiResponseRef: ServiceGradesDeleteResponse
 * @apiSummary: Delete servicegrade
 * @apiDescription: Delete servicegrade
 * @apiSecurity: apiKeySecurity
 */
func (h *ServiceGradeHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ServiceGradesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.ServiceGradeService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
