package handlers

import (
	"errors"
	cPorts "github.com/hoitek/Maja-Service/internal/company/ports"
	ud "github.com/hoitek/Maja-Service/internal/user/domain"
	uModels "github.com/hoitek/Maja-Service/internal/user/models"
	uPorts "github.com/hoitek/Maja-Service/internal/user/ports"
	"github.com/hoitek/Maja-Service/internal/vehicle/models"
	"github.com/hoitek/Maja-Service/internal/vehicle/ports"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/vehicle/config"
)

type VehicleHandler struct {
	VehicleService ports.VehicleService
	UserService    uPorts.UserService
	CompanyService cPorts.CompanyService
}

func NewVehicleHandler(r *mux.Router, s ports.VehicleService, u uPorts.UserService, c cPorts.CompanyService) (VehicleHandler, error) {
	vehicleHandler := VehicleHandler{
		VehicleService: s,
		UserService:    u,
		CompanyService: c,
	}
	if r == nil {
		return VehicleHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.VehicleConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.VehicleConfig.ApiVersion1).Subrouter()

	rv1.Handle("/vehicles", vehicleHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/vehicles", vehicleHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/vehicles", vehicleHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/vehicles/{id}", vehicleHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/vehicles/csv/download", vehicleHandler.Download()).Methods(http.MethodGet)
	return vehicleHandler, nil
}

/*
* @apiTag: vehicle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: vehicles
* @apiResponseRef: VehiclesQueryResponse
* @apiSummary: Query vehicles
* @apiParametersRef: VehiclesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: VehiclesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: VehiclesQueryNotFoundResponse
 */
func (h *VehicleHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.VehiclesQueryRequestParams{}
		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		vehicles, err := h.VehicleService.Query(queries)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(vehicles)
	})
}

/*
* @apiTag: vehicle
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: vehicles/csv/download
* @apiResponseRef: VehiclesQueryResponse
* @apiSummary: Query vehicles
* @apiParametersRef: VehiclesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: VehiclesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: VehiclesQueryNotFoundResponse
 */
func (h *VehicleHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.VehiclesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		vehicles, err := h.VehicleService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(vehicles)
	})
}

/*
 * @apiTag: vehicle
 * @apiPath: /vehicles
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: VehiclesCreateRequestBody
 * @apiResponseRef: VehiclesCreateResponse
 * @apiSummary: Create vehicle
 * @apiDescription: Create vehicle
 */
func (h *VehicleHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.VehiclesCreateRequestBody{}
		errs := payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if user exists
		users, err := h.UserService.Query(&uModels.UsersQueryRequestParams{})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if users.TotalRows == 0 {
			return response.ErrorBadRequest(nil, "user not found")
		}
		payload.User = users.Items.([]*ud.User)[0]

		// Create vehicle
		vehicle, err := h.VehicleService.Create(payload)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(vehicle)
	})
}

/*
 * @apiTag: vehicle
 * @apiPath: /vehicles/{id}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: VehiclesUpdateRequestParams
 * @apiRequestRef: VehiclesCreateRequestBody
 * @apiResponseRef: VehiclesCreateResponse
 * @apiSummary: Update vehicle
 * @apiDescription: Update vehicle
 */
func (h *VehicleHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		// Validate params
		p := mux.Vars(r)
		params := &models.VehiclesUpdateRequestParams{}
		errs := params.ValidateParams(p)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Validate payload
		payload := &models.VehiclesCreateRequestBody{}
		errs = payload.ValidateBody(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		// Check if vehicle exists
		vehicles, err := h.VehicleService.Query(&models.VehiclesQueryRequestParams{
			ID: params.ID,
		})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if vehicles.TotalRows == 0 {
			return response.ErrorNotFound("vehicle not found")
		}

		// Check if user exists
		users, err := h.UserService.Query(&uModels.UsersQueryRequestParams{})
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}
		if users.TotalRows == 0 {
			return response.ErrorBadRequest(nil, "user not found")
		}
		payload.User = users.Items.([]*ud.User)[0]

		// Update vehicle
		data, err := h.VehicleService.Update(payload, params.ID)
		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: vehicle
 * @apiPath: /vehicles
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: VehiclesDeleteRequestBody
 * @apiResponseRef: VehiclesCreateResponse
 * @apiSummary: Delete vehicle
 * @apiDescription: Delete vehicle
 */
func (h *VehicleHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.VehiclesDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.VehicleService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
