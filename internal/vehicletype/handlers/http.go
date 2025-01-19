package handlers

import (
	"errors"
	"github.com/hoitek/Maja-Service/internal/vehicletype/models"
	"github.com/hoitek/Maja-Service/internal/vehicletype/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoitek/Kit/response"
	"github.com/hoitek/Maja-Service/internal/vehicletype/config"
)

type VehicleTypeHandler struct {
	VehicleTypeService service.VehicleTypeService
}

func NewVehicleTypeHandler(r *mux.Router, s service.VehicleTypeService) (VehicleTypeHandler, error) {
	vehicletypeHandler := VehicleTypeHandler{
		VehicleTypeService: s,
	}
	if r == nil {
		return VehicleTypeHandler{}, errors.New("router can not be nil")
	}

	// Leading slash(/) is required for PathPrefix
	rapi := r.PathPrefix(config.VehicleTypeConfig.ApiPrefix).Subrouter()
	rv1 := rapi.PathPrefix(config.VehicleTypeConfig.ApiVersion1).Subrouter()

	rv1.Handle("/vehicletypes", vehicletypeHandler.Create()).Methods(http.MethodPost)
	rv1.Handle("/vehicletypes", vehicletypeHandler.Query()).Methods(http.MethodGet)
	rv1.Handle("/vehicletypes", vehicletypeHandler.Delete()).Methods(http.MethodDelete)
	rv1.Handle("/vehicletypes/{name}", vehicletypeHandler.Update()).Methods(http.MethodPut)
	rv1.Handle("/vehicletypes/csv/download", vehicletypeHandler.Download()).Methods(http.MethodGet)
	return vehicletypeHandler, nil
}

/*
* @apiTag: vehicletype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: vehicletypes
* @apiResponseRef: VehicleTypesQueryResponse
* @apiSummary: Query vehicletypes
* @apiParametersRef: VehicleTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: VehicleTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: VehicleTypesQueryNotFoundResponse
 */
func (h *VehicleTypeHandler) Query() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.VehicleTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		vehicletypes, err := h.VehicleTypeService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(vehicletypes)
	})
}

/*
* @apiTag: vehicletype
* @apiMethod: GET
* @apiStatusCode: 200
* @apiDeprecated: false
* @apiPath: vehicletypes/csv/download
* @apiResponseRef: VehicleTypesQueryResponse
* @apiSummary: Query vehicletypes
* @apiParametersRef: VehicleTypesQueryRequestParams
* @apiErrorStatusCodes: 400, 500, 404
* @api404ResponseRef: VehicleTypesQueryNotFoundResponse
* @api404ResponseDescription: lorem ipsum dolor sit amet
* @api500ResponseRef: VehicleTypesQueryNotFoundResponse
 */
func (h *VehicleTypeHandler) Download() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		queries := &models.VehicleTypesQueryRequestParams{}

		errs := queries.ValidateQueries(r)
		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		vehicletypes, err := h.VehicleTypeService.Query(queries)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(vehicletypes)
	})
}

/*
 * @apiTag: vehicletype
 * @apiPath: /vehicletypes
 * @apiMethod: POST
 * @apiStatusCode: 201
 * @apiRequestRef: VehicleTypesCreateRequestBody
 * @apiResponseRef: VehicleTypesCreateResponse
 * @apiSummary: Create vehicletype
 * @apiDescription: Create vehicletype
 */
func (h *VehicleTypeHandler) Create() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.VehicleTypesCreateRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		user, err := h.VehicleTypeService.Create(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Created(user)
	})
}

/*
 * @apiTag: vehicletype
 * @apiPath: /vehicletypes/{name}
 * @apiMethod: PUT
 * @apiStatusCode: 201
 * @apiParametersRef: VehicleTypesUpdateRequestParams
 * @apiRequestRef: VehicleTypesCreateRequestBody
 * @apiResponseRef: VehicleTypesCreateResponse
 * @apiSummary: Update vehicletype
 * @apiDescription: Update vehicletype
 */
func (h *VehicleTypeHandler) Update() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		p := mux.Vars(r)
		params := &models.VehicleTypesUpdateRequestParams{}
		errs := params.ValidateParams(p)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		payload := &models.VehicleTypesCreateRequestBody{}
		errs = payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.VehicleTypeService.Update(payload, params.Name)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}

/*
 * @apiTag: vehicletype
 * @apiPath: /vehicletypes
 * @apiMethod: DELETE
 * @apiStatusCode: 201
 * @apiRequestRef: VehicleTypesDeleteRequestBody
 * @apiResponseRef: VehicleTypesCreateResponse
 * @apiSummary: Delete vehicletype
 * @apiDescription: Delete vehicletype
 */
func (h *VehicleTypeHandler) Delete() response.Handler {
	return response.Handle(func(w http.ResponseWriter, r *http.Request) response.Response {
		payload := &models.VehicleTypesDeleteRequestBody{}
		errs := payload.ValidateBody(r)

		if len(errs) > 0 {
			return response.ErrorBadRequest(errs, "Your request data is invalid")
		}

		data, err := h.VehicleTypeService.Delete(payload)

		if err != nil {
			return response.ErrorInternalServerError(nil, err.Error())
		}

		return response.Success(data)
	})
}
