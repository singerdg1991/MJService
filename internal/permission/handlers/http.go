package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/permission/models"
	"github.com/hoitek/Maja-Service/internal/permission/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/permission/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type PermissionHandler struct {
	UserService       uPorts.UserService
	PermissionService ports.PermissionService
}

func NewPermissionHandler(r *mux.Router, s ports.PermissionService, us uPorts.UserService) (PermissionHandler, error) {
	permissionHandler := PermissionHandler{
		UserService:       us,
		PermissionService: s,
	}
	if r == nil {
		return PermissionHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.PermissionConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.PermissionConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/permissions",
			Method:      http.MethodPost,
			Handler:     permissionHandler.Create(),
			Permissions: []string{permissions.STAFF_PERMISSIONS_CREATE_NEW_STAFF_PERMISSION},
		},
		{
			Path:        "/permissions",
			Method:      http.MethodGet,
			Handler:     permissionHandler.Query(),
			Permissions: []string{permissions.STAFF_PERMISSIONS_VIEW_ALL_STAFF_PERMISSIONS},
		},
		{
			Path:        "/permissions",
			Method:      http.MethodDelete,
			Handler:     permissionHandler.Delete(),
			Permissions: []string{permissions.STAFF_PERMISSIONS_CREATE_NEW_STAFF_PERMISSION},
		},
		{
			Path:        "/permissions/{id}",
			Method:      http.MethodPut,
			Handler:     permissionHandler.Update(),
			Permissions: []string{permissions.STAFF_PERMISSIONS_CREATE_NEW_STAFF_PERMISSION},
		},
		{
			Path:        "/permissions/csv/download",
			Method:      http.MethodGet,
			Handler:     permissionHandler.Download(),
			Permissions: []string{permissions.STAFF_PERMISSIONS_VIEW_ALL_STAFF_PERMISSIONS},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(us, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return permissionHandler, nil
}

/*
* @apiTag: permission
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: permissions
* @apiResponseRef: PermissionsQueryResponse
* @apiSummary: Query permissions
* @apiParametersRef: PermissionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PermissionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PermissionsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *PermissionHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PermissionsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query permissions
		permissions, err := h.PermissionService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(permissions)
	})
}

/*
* @apiTag: permission
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: permissions/csv/download
* @apiResponseRef: PermissionsQueryResponse
* @apiSummary: Query permissions
* @apiParametersRef: PermissionsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: PermissionsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: PermissionsQueryNotFoundResponse
 */
func (h *PermissionHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.PermissionsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		permissions, err := h.PermissionService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(permissions)
	})
}

/*
 * @apiTag: permission
 * @apiPath: /permissions
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: PermissionsCreateRequestBody
 * @apiResponseRef: PermissionsCreateResponse
 * @apiSummary: Create permission
 * @apiDescription: Create permission
 * @apiSecurity: apiKeySecurity
 */
func (h *PermissionHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PermissionsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find permission by name
		res, _ := h.PermissionService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "Permission already exists")
		}

		// Create permission
		permission, err := h.PermissionService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(permission)
	})
}

/*
 * @apiTag: permission
 * @apiPath: /permissions/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: PermissionsUpdateRequestParams
 * @apiRequestRef: PermissionsCreateRequestBody
 * @apiResponseRef: PermissionsCreateResponse
 * @apiSummary: Update permission
 * @apiDescription: Update permission
 * @apiSecurity: apiKeySecurity
 */
func (h *PermissionHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.PermissionsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.PermissionsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.PermissionService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.PermissionService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.PermissionService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: permission
 * @apiPath: /permissions
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: PermissionsDeleteRequestBody
 * @apiResponseRef: PermissionsDeleteResponse
 * @apiSummary: Delete permission
 * @apiDescription: Delete permission
 * @apiSecurity: apiKeySecurity
 */
func (h *PermissionHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.PermissionsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.PermissionService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
