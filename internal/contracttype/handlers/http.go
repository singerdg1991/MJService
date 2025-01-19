package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/contracttype/models"
	"github.com/hoitek/Maja-Service/internal/contracttype/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/contracttype/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type ContractTypeHandler struct {
	ContractTypeService ports.ContractTypeService
	UserService         uPorts.UserService
}

func NewContractTypeHandler(r *mux.Router, s ports.ContractTypeService, u uPorts.UserService) (ContractTypeHandler, error) {
	contractTypeHandler := ContractTypeHandler{
		ContractTypeService: s,
		UserService:         u,
	}
	if r == nil {
		return ContractTypeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.ContractTypeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.ContractTypeConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/contracttypes",
			Method:      http.MethodPost,
			Handler:     contractTypeHandler.Create(),
			Permissions: []string{permissions.CONTRACT_TYPES_CREATE_NEW_CONTRACT_TYPE},
		},
		{
			Path:        "/contracttypes",
			Method:      http.MethodGet,
			Handler:     contractTypeHandler.Query(),
			Permissions: []string{permissions.CONTRACT_TYPES_VIEW_ALL_CONTRACT_TYPES},
		},
		{
			Path:        "/contracttypes",
			Method:      http.MethodDelete,
			Handler:     contractTypeHandler.Delete(),
			Permissions: []string{permissions.CONTRACT_TYPES_CREATE_NEW_CONTRACT_TYPE},
		},
		{
			Path:        "/contracttypes/{id}",
			Method:      http.MethodPut,
			Handler:     contractTypeHandler.Update(),
			Permissions: []string{permissions.CONTRACT_TYPES_CREATE_NEW_CONTRACT_TYPE},
		},
		{
			Path:        "/contracttypes/csv/download",
			Method:      http.MethodGet,
			Handler:     contractTypeHandler.Download(),
			Permissions: []string{permissions.CONTRACT_TYPES_VIEW_ALL_CONTRACT_TYPES},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return contractTypeHandler, nil
}

/*
* @apiTag: contracttype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: contracttypes
* @apiResponseRef: ContractTypesQueryResponse
* @apiSummary: Query ContractTypes
* @apiParametersRef: ContractTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ContractTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ContractTypesQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *ContractTypeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ContractTypesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query ContractTypes
		contractTypes, err := h.ContractTypeService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(contractTypes)
	})
}

/*
* @apiTag: contracttype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: contracttypes/csv/download
* @apiResponseRef: ContractTypesQueryResponse
* @apiSummary: Query ContractTypes
* @apiParametersRef: ContractTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: ContractTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: ContractTypesQueryNotFoundResponse
 */
func (h *ContractTypeHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.ContractTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		contractTypes, err := h.ContractTypeService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(contractTypes)
	})
}

/*
 * @apiTag: contracttype
 * @apiPath: /contracttypes
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: ContractTypesCreateRequestBody
 * @apiResponseRef: ContractTypesCreateResponse
 * @apiSummary: Create ContractTypes
 * @apiDescription: Create ContractTypes
 * @apiSecurity: apiKeySecurity
 */
func (h *ContractTypeHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ContractTypesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find ContractTypes by name
		res, _ := h.ContractTypeService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "ContractType already exists")
		}

		// Create ContractTypes
		contractTypes, err := h.ContractTypeService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(contractTypes)
	})
}

/*
 * @apiTag: contracttype
 * @apiPath: /contracttypes/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: ContractTypesUpdateRequestParams
 * @apiRequestRef: ContractTypesCreateRequestBody
 * @apiResponseRef: ContractTypesCreateResponse
 * @apiSummary: Update ContractTypes
 * @apiDescription: Update ContractTypes
 * @apiSecurity: apiKeySecurity
 */
func (h *ContractTypeHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.ContractTypesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.ContractTypesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.ContractTypeService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.ContractTypeService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.ContractTypeService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: contracttype
 * @apiPath: /contracttypes
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: ContractTypesDeleteRequestBody
 * @apiResponseRef: ContractTypesDeleteResponse
 * @apiSummary: Delete ContractTypes
 * @apiDescription: Delete ContractTypes
 * @apiSecurity: apiKeySecurity
 */
func (h *ContractTypeHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.ContractTypesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.ContractTypeService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
