package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hoitek/Maja-Service/internal/_shared/route"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/permissions"

	"github.com/hoitek/Maja-Service/internal/_shared/middlewares"
	"github.com/hoitek/Maja-Service/internal/equipment/models"
	"github.com/hoitek/Maja-Service/internal/equipment/ports"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/equipment/config"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
)

type EquipmentHandler struct {
	EquipmentService ports.EquipmentService
	UserService      uPorts.UserService
}

func NewEquipmentHandler(r *mux.Router, s ports.EquipmentService, u uPorts.UserService) (EquipmentHandler, error) {
	equipmentHandler := EquipmentHandler{
		EquipmentService: s,
		UserService:      u,
	}
	if r == nil {
		return EquipmentHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.EquipmentConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.EquipmentConfig.ApiVersion1).Subrouter()

	// Create secure routes
	secureRoutes := []route.SecureRoute{
		{
			Path:        "/equipments",
			Method:      http.MethodPost,
			Handler:     equipmentHandler.Create(),
			Permissions: []string{permissions.EQUIPMENTS_CREATE_NEW_EQUIPMENT},
		},
		{
			Path:        "/equipments",
			Method:      http.MethodGet,
			Handler:     equipmentHandler.Query(),
			Permissions: []string{permissions.EQUIPMENTS_VIEW_ALL_EQUIPMENTS},
		},
		{
			Path:        "/equipments",
			Method:      http.MethodDelete,
			Handler:     equipmentHandler.Delete(),
			Permissions: []string{permissions.EQUIPMENTS_CREATE_NEW_EQUIPMENT},
		},
		{
			Path:        "/equipments/{id}",
			Method:      http.MethodPut,
			Handler:     equipmentHandler.Update(),
			Permissions: []string{permissions.EQUIPMENTS_CREATE_NEW_EQUIPMENT},
		},
		{
			Path:        "/equipments/csv/download",
			Method:      http.MethodGet,
			Handler:     equipmentHandler.Download(),
			Permissions: []string{permissions.EQUIPMENTS_VIEW_ALL_EQUIPMENTS},
		},
	}

	// Register secure routes
	for _, route := range secureRoutes {
		rAuth := rv1.Path(route.Path).Handler(middlewares.OAuth2Middleware(middlewares.AuthMiddleware(u, route.Permissions)(route.Handler)))
		rAuth.Methods(route.Method)
	}

	return equipmentHandler, nil
}

/*
* @apiTag: equipment
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: equipments
* @apiResponseRef: EquipmentsQueryResponse
* @apiSummary: Query equipments
* @apiParametersRef: EquipmentsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: EquipmentsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: EquipmentsQueryNotFoundResponse
* @apiSecurity: apiKeySecurity
 */
func (h *EquipmentHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.EquipmentsQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Query equipments
		equipments, err := h.EquipmentService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(equipments)
	})
}

/*
* @apiTag: equipment
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: equipments/csv/download
* @apiResponseRef: EquipmentsQueryResponse
* @apiSummary: Query equipments
* @apiParametersRef: EquipmentsQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: EquipmentsQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: EquipmentsQueryNotFoundResponse
 */
func (h *EquipmentHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.EquipmentsQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		equipments, err := h.EquipmentService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(equipments)
	})
}

/*
 * @apiTag: equipment
 * @apiPath: /equipments
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: EquipmentsCreateRequestBody
 * @apiResponseRef: EquipmentsCreateResponse
 * @apiSummary: Create equipment
 * @apiDescription: Create equipment
 * @apiSecurity: apiKeySecurity
 */
func (h *EquipmentHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.EquipmentsCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Find equipment by name
		res, _ := h.EquipmentService.FindByName(payload.Name)
		if res != nil {
			return response.ErrorInternalServerError(nil, "Equipment already exists")
		}

		// Create equipment
		equipment, err := h.EquipmentService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(equipment)
	})
}

/*
 * @apiTag: equipment
 * @apiPath: /equipments/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: EquipmentsUpdateRequestParams
 * @apiRequestRef: EquipmentsCreateRequestBody
 * @apiResponseRef: EquipmentsCreateResponse
 * @apiSummary: Update equipment
 * @apiDescription: Update equipment
 * @apiSecurity: apiKeySecurity
 */
func (h *EquipmentHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.EquipmentsUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.EquipmentsCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if Language Skill already exists
		res, _ := h.EquipmentService.FindByID(int64(params.ID))
		if res == nil {
			return response.ErrorBadRequest(nil, "Language Skill does not exist")
		}

		// Check if Language Skill already exists with the same name
		res, _ = h.EquipmentService.FindByName(payload.Name)
		log.Println(res, payload.Name, params.ID)
		if res != nil && int64(res.ID) != int64(params.ID) {
			return response.ErrorBadRequest(nil, "Language Skill already exists")
		}

		// Update Language Skill
		data, err := h.EquipmentService.Update(payload, int64(params.ID))
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: equipment
 * @apiPath: /equipments
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: EquipmentsDeleteRequestBody
 * @apiResponseRef: EquipmentsDeleteResponse
 * @apiSummary: Delete equipment
 * @apiDescription: Delete equipment
 * @apiSecurity: apiKeySecurity
 */
func (h *EquipmentHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.EquipmentsDeleteRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		ids, err := utils.ConvertInterfaceSliceToSliceOfInt64(payload.IDs)
		if err != nil {
			return response.ErrorBadRequest(nil, "IDs is invalid")
		}
		payload.IDsInt64 = ids

		data, err := h.EquipmentService.Delete(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
