package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	permPorts "github.com/hoitek/Maja-Service/internal/permission/ports"
	"github.com/hoitek/Maja-Service/internal/role/models"
	rPorts "github.com/hoitek/Maja-Service/internal/role/ports"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/role/config"
)

type RoleHandler struct {
	RoleService       rPorts.RoleService
	PermissionService permPorts.PermissionService
}

func NewRoleHandler(r *mux.Router, s rPorts.RoleService, p permPorts.PermissionService) (RoleHandler, error) {
	roleHandler := RoleHandler{
		RoleService:       s,
		PermissionService: p,
	}
	if r == nil {
		return RoleHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.RoleConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.RoleConfig.ApiVersion1).Subrouter()

	rv1.Handle("/roles", roleHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/roles", roleHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/roles", roleHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/roles/{id}", roleHandler.Update()).Methods(http.MethodPut)
	return roleHandler, nil
}

/*
* @apiTag: role
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: roles
* @apiResponseRef: RolesQueryResponse
* @apiSummary: Query roles
* @apiParametersRef: RolesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: RolesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: RolesQueryNotFoundResponse
 */
func (h *RoleHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.RolesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		roles, err := h.RoleService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(roles)
	})
}

/*
 * @apiTag: role
 * @apiPath: /roles
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: RolesCreateRequestBody
 * @apiResponseRef: RolesCreateResponse
 * @apiSummary: Create role
 * @apiDescription: Create role
 */
func (h *RoleHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.RolesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check role exists
		foundRole := h.RoleService.GetRoleByName(payload.Name)
		if foundRole != nil {
			return response.ErrorNotFound(nil, "Role already exists")
		}

		// Convert interface slice to slice of int64
		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.Permissions)
		if err != nil {
			return response.ErrorBadRequest(nil, "Permissions is invalid")
		}
		payload.PermissionsInt64 = ids

		// Check Permissions Exists
		permissions, err := h.PermissionService.GetPermissionsByIds(payload.PermissionsInt64)
		if err != nil {
			log.Printf("Error when get permissions by ids: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}
		if len(permissions) != len(payload.PermissionsInt64) {
			return response.ErrorBadRequest(nil, "Permissions is invalid")
		}

		// Create role
		insertedRole, err := h.RoleService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(insertedRole)
	})
}

/*
 * @apiTag: role
 * @apiPath: /roles/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 200
 * @apiParametersRef: RolesUpdateRequestParams
 * @apiRequestRef: RolesCreateRequestBody
 * @apiResponseRef: RolesCreateResponse
 * @apiSummary: Update role
 * @apiDescription: Update role
 */
func (h *RoleHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.RolesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate body
		payload := &models.RolesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Convert interface slice to slice of int64
		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.Permissions)
		if err != nil {
			return response.ErrorBadRequest(nil, "Permissions is invalid")
		}
		payload.PermissionsInt64 = ids

		// Check Permissions Exists
		permissions, err := h.PermissionService.GetPermissionsByIds(payload.PermissionsInt64)
		if err != nil {
			log.Printf("Error when get permissions by ids: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}
		if len(permissions) != len(payload.PermissionsInt64) {
			return response.ErrorBadRequest(nil, "Permissions is invalid")
		}

		// Check role exists
		foundRole := h.RoleService.GetRoleByID(params.ID)
		if foundRole == nil {
			return response.ErrorNotFound(nil, "Role not found")
		}

		// Validate role name
		foundRoleByName := h.RoleService.GetRoleByName(payload.Name)
		if foundRoleByName != nil && foundRoleByName.ID != foundRole.ID {
			return response.ErrorBadRequest(nil, "Role name is already exists")
		}

		// Update role
		data, err := h.RoleService.Update(payload, int64(params.ID))
		if err != nil {
			log.Printf("Error when update role: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: role
 * @apiPath: /roles
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: RolesDeleteRequestBody
 * @apiResponseRef: RolesCreateResponse
 * @apiSummary: Delete role
 * @apiDescription: Delete role
 */
func (h *RoleHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.RolesDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Convert interface slice to slice of int64
		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "Permissions is invalid")
		}
		payload.IDsInt64 = ids

		// Delete role
		data, err := h.RoleService.Delete(payload)
		if err != nil {
			log.Printf("Error when delete role: %v\n", err)
			return response.ErrorInternalServerError(nil, "Something went wrong, please try again later")
		}

		return response.Success(data)
	})
}
